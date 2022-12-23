package main
// This template draws inspiration from the following source:
// github.com/lucasfth/go-ass4
// Which was created by chbl, fefa and luha

import (
	request "Lucasfth/disys-exam-peer-to-peer/grpc/grpc"
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.Lmicroseconds)

	// Get port to listen on
	log.Print("Enter id of peer (1-3), below:")
	var ownPort int32
	fmt.Scanln(&ownPort)
	ownPort = 5000 + ownPort

	// To change log location, outcomment below

	// path := fmt.Sprintf("clientlog_%v", ownPort)
	// f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()

	// log.SetOutput(f)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		port:          ownPort,
		requestAmount: 0,
		isPiloting:    false,
		peers:         make(map[int32]request.RequestServiceClient),
		ctx:           ctx,
	}

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", p.port))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	grpcServer := grpc.NewServer()
	request.RegisterRequestServiceServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	p.connectToPeers() // Will not finish until all peers are connected

	p.interactWithPeers()
}

func (p *peer) interactWithPeers() {
	for {
		responses, shouldTry := p.sendRequestToAll()

		if !shouldTry {
			continue
		}

		for i := 0; i < 4; i++ {
			if i == 3 {
				p.criticalSection()
				p.requestAmount = 0
				break
			}

			if int32(i)+5001 == p.port {
				continue
			}
			if responses[int32(i)+5001] > p.requestAmount {
				p.mutex.Unlock()
				log.Printf("unlocked")
				break
			} else if responses[int32(i)+5001] == p.requestAmount && int32(i)+5001 > p.port {
				p.mutex.Unlock()
				log.Printf("unlocked")
				break
			}
		}
	}
}

func (p *peer) criticalSection() {
	p.isPiloting = true
	log.Printf("%v is now pilot 	-----------------------", p.port)
	time.Sleep(4 * time.Second)
	log.Printf("%v is not pilot 	-----------", p.port)
	p.isPiloting = false
	p.mutex.Unlock()
	log.Printf("unlocked")
	time.Sleep(2 * time.Second)
}

func (p *peer) connectToPeers() {
	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i+1)

		if port == p.port {
			continue
		}

		var conn *grpc.ClientConn
		log.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := request.NewRequestServiceClient(conn)
		p.peers[port] = c
	}
}

func (p *peer) Request(ctx context.Context, req *request.Request) (*request.Reply, error) {
	id := req.Id
	reqAmount := p.requestAmount

	rep := &request.Reply{Id: id, RequestAmount: reqAmount, IsPiloting: p.isPiloting}
	return rep, nil
}

func (p *peer) sendRequestToAll() (map[int32]int32, bool) {
	response := make(map[int32]int32)
	p.requestAmount++
	
	p.mutex.Lock()
	log.Printf("locked")

	request := &request.Request{Id: p.port, RequestAmount: p.requestAmount}
	for id, peer := range p.peers {
		reply, err := peer.Request(p.ctx, request)
		if err != nil {
			log.Fatalf("Could not send request: %v", err)
		}
		log.Printf("Got reply from id %v: %v: %v\n", id, reply.RequestAmount, reply.IsPiloting)
		if reply.IsPiloting {
			p.mutex.Unlock()
			log.Printf("unlock")
			time.Sleep(2 * time.Second)
			return make(map[int32]int32), false
		}
		response[reply.Id] = reply.RequestAmount
	}
	return response, true
}

type peer struct {
	request.UnimplementedRequestServiceServer
	mutex         sync.Mutex
	port          int32
	requestAmount int32
	isPiloting    bool
	peers         map[int32]request.RequestServiceClient
	ctx           context.Context
}
