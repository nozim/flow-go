package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dapperlabs/bamboo-emulator/crypto"
	"github.com/dapperlabs/bamboo-emulator/data"
	"github.com/dapperlabs/bamboo-emulator/gen/grpc/services/accessv1"
	"github.com/dapperlabs/bamboo-emulator/nodes/access"
)

// Server is a gRPC server that implements the Bamboo Access API.
type Server interface {
	// Start runs the server at the provided port.
	Start(port int)
}

type server struct {
	accessNode access.Node
}

// NewServer returns a new Bamboo emulator server.
func NewServer(accessNode access.Node) Server {
	return &server{
		accessNode: accessNode,
	}
}

// SendTransaction submits a transaction to the network.
func (s *server) SendTransaction(ctx context.Context, req *accessv1.SendTransactionRequest) (*accessv1.SendTransactionResponse, error) {
	txMsg := req.GetTransaction()

	tx := &data.Transaction{
		ToAddress:      crypto.BytesToAddress(txMsg.GetTo()),
		Script:         txMsg.GetScript(),
		Nonce:          txMsg.GetNonce(),
		ComputeLimit:   txMsg.GetComputeLimit(),
		PayerSignature: txMsg.GetPayerSignature(),
		Status:         data.TxPending,
	}

	err := s.accessNode.SendTransaction(tx)
	if err != nil {
		switch err.(type) {
		case *access.DuplicateTransactionError:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &accessv1.SendTransactionResponse{
		Hash: tx.Hash().Bytes(),
	}, nil
}

// GetBlockByHash gets a block by hash.
func (s *server) GetBlockByHash(ctx context.Context, req *accessv1.GetBlockByHashRequest) (*accessv1.GetBlockByHashResponse, error) {
	hash := crypto.BytesToHash(req.GetHash())

	block, err := s.accessNode.GetBlockByHash(hash)
	if err != nil {
		switch err.(type) {
		case *access.BlockNotFoundError:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &accessv1.GetBlockByHashResponse{
		Block: &accessv1.Block{
			Hash:              block.Hash().Bytes(),
			Number:            block.Number,
			TransactionHashes: crypto.HashesToBytes(block.TransactionHashes),
			Status:            accessv1.Block_Status(block.Status),
		},
	}, nil
}

// GetBlockByNumber gets a block by number.
func (s *server) GetBlockByNumber(ctx context.Context, req *accessv1.GetBlockByNumberRequest) (*accessv1.GetBlockByNumberResponse, error) {
	number := req.GetNumber()

	block, err := s.accessNode.GetBlockByNumber(number)
	if err != nil {
		switch err.(type) {
		case *access.BlockNotFoundError:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &accessv1.GetBlockByNumberResponse{
		Block: &accessv1.Block{
			Hash:              block.Hash().Bytes(),
			Number:            block.Number,
			TransactionHashes: crypto.HashesToBytes(block.TransactionHashes),
			Status:            accessv1.Block_Status(block.Status),
		},
	}, nil
}

// GetLatestBlock gets the latest sealed block.
func (s *server) GetLatestBlock(ctx context.Context, req *accessv1.GetLatestBlockRequest) (*accessv1.GetLatestBlockResponse, error) {
	block := s.accessNode.GetLatestBlock()

	return &accessv1.GetLatestBlockResponse{
		Block: &accessv1.Block{
			Hash:              block.Hash().Bytes(),
			Number:            block.Number,
			TransactionHashes: crypto.HashesToBytes(block.TransactionHashes),
			Status:            accessv1.Block_Status(block.Status),
		},
	}, nil
}

// GetTransactions gets a transaction by hash.
func (s *server) GetTransaction(ctx context.Context, req *accessv1.GetTransactionRequest) (*accessv1.GetTransactionResponse, error) {
	hash := crypto.BytesToHash(req.GetHash())

	tx, err := s.accessNode.GetTransaction(hash)
	if err != nil {
		switch err.(type) {
		case *access.TransactionNotFoundError:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &accessv1.GetTransactionResponse{
		Transaction: &accessv1.GetTransactionResponse_Transaction{
			To:             tx.ToAddress.Bytes(),
			Script:         tx.Script,
			Nonce:          tx.Nonce,
			ComputeLimit:   tx.ComputeLimit,
			ComputeUsed:    tx.ComputeUsed,
			PayerSignature: tx.PayerSignature,
			Status:         accessv1.GetTransactionResponse_Transaction_Status(tx.Status),
		},
	}, nil
}

// GetBalance returns the balance of an address.
func (s *server) GetBalance(context.Context, *accessv1.GetBalanceRequest) (*accessv1.GetBalanceResponse, error) {
	var address crypto.Address

	s.accessNode.GetBalance(address)

	return nil, nil
}

// CallContract performs a contract call.
func (s *server) CallContract(context.Context, *accessv1.CallContractRequest) (*accessv1.CallContractResponse, error) {
	return nil, nil
}

func (s *server) Start(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	accessv1.RegisterBambooAccessAPIServer(grpcServer, s)

	grpcServer.Serve(lis)
}
