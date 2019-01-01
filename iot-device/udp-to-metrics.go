package main

func main() {
	inbound, inpktpool, err := Listen(":7331", nil)
	if err != nil {
		// handle err
	}

	// Receive UDP packet
	inpkt := <-inbound

	// Do something with UDP packet

	// Release UDP packet back to pool
	inpktpool.Release(inpkt)

	return
}
