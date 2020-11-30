package blockcreator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.ibm.com/blockchaindb/server/internal/blockstore"
	"github.ibm.com/blockchaindb/server/internal/queue"
	"github.ibm.com/blockchaindb/server/internal/worldstate"
	"github.ibm.com/blockchaindb/server/internal/worldstate/leveldb"
	"github.ibm.com/blockchaindb/server/pkg/logger"
	"github.ibm.com/blockchaindb/server/pkg/types"
)

type testEnv struct {
	creator        *BlockCreator
	db             worldstate.DB
	dbPath         string
	blockStore     *blockstore.Store
	blockStorePath string
	cleanup        func()
}

func newTestEnv(t *testing.T) *testEnv {
	c := &logger.Config{
		Level:         "debug",
		OutputPath:    []string{"stdout"},
		ErrOutputPath: []string{"stderr"},
		Encoding:      "console",
	}
	logger, err := logger.New(c)
	require.NoError(t, err)

	dir, err := ioutil.TempDir("/tmp", "creator")
	require.NoError(t, err)

	dbPath := filepath.Join(dir, "leveldb")
	db, err := leveldb.Open(
		&leveldb.Config{
			DBRootDir: dbPath,
			Logger:    logger,
		},
	)
	if err != nil {
		if rmErr := os.RemoveAll(dir); rmErr != nil {
			t.Errorf("error while removing directory %s, %v", dir, err)
		}
		t.Fatalf("error while creating the leveldb instance, %v", err)
	}

	blockStorePath := filepath.Join(dir, "blockstore")
	blockStore, err := blockstore.Open(
		&blockstore.Config{
			StoreDir: blockStorePath,
			Logger:   logger,
		},
	)
	if err != nil {
		if rmErr := os.RemoveAll(dir); rmErr != nil {
			t.Errorf("error while removing directory %s, %v", dir, err)
		}
		t.Fatalf("error while creating the block store, %v", err)
	}

	b, err := New(&Config{
		TxBatchQueue: queue.New(10),
		BlockQueue:   queue.New(10),
		Logger:       logger,
		BlockStore:   blockStore,
	})
	require.NoError(t, err)
	go b.Start()
	b.WaitTillStart()

	cleanup := func() {
		b.Stop()

		if err := db.Close(); err != nil {
			t.Errorf("failed to close the db instance, %v", err)
		}

		if err := blockStore.Close(); err != nil {
			t.Errorf("failed to close the blockstore, %v", err)
		}

		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("failed to remove directory %s, %v", dir, err)
		}
	}

	return &testEnv{
		creator:        b,
		db:             db,
		dbPath:         dir,
		blockStore:     blockStore,
		blockStorePath: blockStorePath,
		cleanup:        cleanup,
	}
}

func TestBatchCreator(t *testing.T) {
	dataTx1 := &types.DataTxEnvelope{
		Payload: &types.DataTx{
			UserID: "user1",
			DBName: "db1",
			DataDeletes: []*types.DataDelete{
				{
					Key: "key1",
				},
			},
		},
	}

	dataTx2 := &types.DataTxEnvelope{
		Payload: &types.DataTx{
			UserID: "user2",
			DBName: "db2",
			DataDeletes: []*types.DataDelete{
				{
					Key: "key2",
				},
			},
		},
	}

	userAdminTx := &types.UserAdministrationTxEnvelope{
		Payload: &types.UserAdministrationTx{
			UserID: "user1",
			UserReads: []*types.UserRead{
				{
					UserID: "user1",
				},
			},
			UserWrites: []*types.UserWrite{
				{
					User: &types.User{
						ID:          "user2",
						Certificate: []byte("certificate"),
					},
				},
			},
		},
	}

	dbAdminTx := &types.DBAdministrationTxEnvelope{
		Payload: &types.DBAdministrationTx{
			UserID:    "user1",
			CreateDBs: []string{"db1", "db2"},
			DeleteDBs: []string{"db3", "db4"},
		},
	}

	configTx := &types.ConfigTxEnvelope{
		Payload: &types.ConfigTx{
			UserID: "user1",
			NewConfig: &types.ClusterConfig{
				Nodes: []*types.NodeConfig{
					{
						ID: "node1",
					},
				},
				Admins: []*types.Admin{
					{
						ID: "admin1",
					},
				},
				RootCACertificate: []byte("root-ca"),
			},
		},
	}

	testCases := []struct {
		name           string
		txBatches      []interface{}
		expectedBlocks []*types.Block
	}{
		{
			name: "enqueue all types of transactions",
			txBatches: []interface{}{
				&types.Block_UserAdministrationTxEnvelope{
					UserAdministrationTxEnvelope: userAdminTx,
				},
				&types.Block_DBAdministrationTxEnvelope{
					DBAdministrationTxEnvelope: dbAdminTx,
				},
				&types.Block_DataTxEnvelopes{
					DataTxEnvelopes: &types.DataTxEnvelopes{
						Envelopes: []*types.DataTxEnvelope{
							dataTx1,
							dataTx2,
						},
					},
				},
				&types.Block_ConfigTxEnvelope{
					ConfigTxEnvelope: configTx,
				},
			},
			expectedBlocks: []*types.Block{
				{
					Header: &types.BlockHeader{
						BaseHeader: &types.BlockHeaderBase{
							Number:                1,
							LastCommittedBlockNum: 0,
						},
						ValidationInfo: []*types.ValidationInfo{
							{
								Flag: types.Flag_VALID,
							},
						},
					},
					Payload: &types.Block_UserAdministrationTxEnvelope{
						UserAdministrationTxEnvelope: userAdminTx,
					},
				},
				{
					Header: &types.BlockHeader{
						BaseHeader: &types.BlockHeaderBase{
							Number: 2,
						},
						ValidationInfo: []*types.ValidationInfo{
							{
								Flag: types.Flag_VALID,
							},
						},
					},
					Payload: &types.Block_DBAdministrationTxEnvelope{
						DBAdministrationTxEnvelope: dbAdminTx,
					},
				},
				{
					Header: &types.BlockHeader{
						BaseHeader: &types.BlockHeaderBase{
							Number: 3,
						},
						ValidationInfo: []*types.ValidationInfo{
							{
								Flag: types.Flag_VALID,
							},
							{
								Flag: types.Flag_VALID,
							},
						},
					},
					Payload: &types.Block_DataTxEnvelopes{
						DataTxEnvelopes: &types.DataTxEnvelopes{
							Envelopes: []*types.DataTxEnvelope{
								dataTx1,
								dataTx2,
							},
						},
					},
				},
				{
					Header: &types.BlockHeader{
						BaseHeader: &types.BlockHeaderBase{
							Number: 4,
						},
						ValidationInfo: []*types.ValidationInfo{
							{
								Flag: types.Flag_VALID,
							},
						},
					},
					Payload: &types.Block_ConfigTxEnvelope{
						ConfigTxEnvelope: configTx,
					},
				},
			},
		},
	}
	for _, tt := range testCases {
		// updating test case blocks PrevCommitted block. We make it block 0, to keep thing simple
		genesisBlockHash, err := blockstore.ComputeBlockHash(tt.expectedBlocks[0])
		require.NoError(t, err)
		for _, expectedBlock := range tt.expectedBlocks[1:] {
			expectedBlock.Header.BaseHeader.LastCommittedBlockNum = 1
			expectedBlock.Header.BaseHeader.LastCommittedBlockHash = genesisBlockHash
		}
		// Update test case blocks prev hashes
		for index, expectedBlock := range tt.expectedBlocks {
			var prevBlockHash []byte
			prevBlockHash = nil
			if index > 0 {
				var err error
				prevBlockHash, err = blockstore.ComputeBlockBaseHash(tt.expectedBlocks[index-1])
				require.NoError(t, err)
			}
			expectedBlock.Header.BaseHeader.PreviousBaseHeaderHash = prevBlockHash
		}
	}

	enqueueTxBatchesAndAssertBlocks := func(t *testing.T, testEnv *testEnv, txBatches []interface{}, expectedBlocks []*types.Block) {
		for _, txBatch := range txBatches {
			testEnv.creator.txBatchQueue.Enqueue(txBatch)
		}

		hasBlockCountMatched := func() bool {
			if len(expectedBlocks) == testEnv.creator.blockQueue.Size() {
				return true
			}
			return false
		}
		require.Eventually(t, hasBlockCountMatched, 2*time.Second, 1000*time.Millisecond)

		for _, expectedBlock := range expectedBlocks {
			block := testEnv.creator.blockQueue.Dequeue().(*types.Block)
			block.Header.ValidationInfo = expectedBlock.Header.ValidationInfo
			require.True(t, proto.Equal(expectedBlock, block), "Expected block  %v, received block %v", expectedBlock, block)
		}
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testEnv := newTestEnv(t)
			defer testEnv.cleanup()

			// storing only first block in block store, to simulate last committed block
			require.NoError(t, testEnv.blockStore.Commit(tt.expectedBlocks[0]))

			enqueueTxBatchesAndAssertBlocks(t, testEnv, tt.txBatches, tt.expectedBlocks)
		})
	}
}
