/*  The MIT License (MIT)

Copyright (c) 2014 Luleå University of Technology, Sweden

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE. */

package soil

import (
	//"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"hash"
	"math/big"
	"testing"
	"time"
)

type Vnode struct {
	Id   []byte // Virtual ID
	Host string // Host identifier
}

// Converts the ID to string
func (vn *Vnode) String() string {
	return fmt.Sprintf("%x", vn.Id)
}

// Generates an ID for the node
func (vn *Vnode) genId(idx uint16, conf *Config) {
	hash := conf.HashFunc()
	hash.Write([]byte(conf.Hostname))
	binary.Write(hash, binary.BigEndian, idx)

	// Use the hash as the ID
	vn.Id = hash.Sum(nil)
}

// Configuration for Chord nodes
type Config struct {
	Hostname      string           // Local host name
	NumVnodes     int              // Number of vnodes per physical node
	HashFunc      func() hash.Hash // Hash function to use
	StabilizeMin  time.Duration    // Minimum stabilization time
	StabilizeMax  time.Duration    // Maximum stabilization time
	NumSuccessors int              // Number of successors to maintain
	Delegate      Delegate         // Invoked to handle ring events
	hashBits      int              // Bit size of the hash function
}

// Delegate to notify on ring events
type Delegate interface {
	/*NewPredecessor(local, remoteNew, remotePrev *Vnode)
	Leaving(local, pred, succ *Vnode)
	PredecessorLeaving(local, remote *Vnode)
	SuccessorLeaving(local, remote *Vnode)
	Shutdown() */
}

// Returns the default Ring configuration
func DefaultConfig(hostname string) *Config {
	return &Config{
		hostname,
		8,        // 8 vnodes
		sha1.New, // SHA1
		time.Duration(15 * time.Second),
		time.Duration(45 * time.Second),
		8,   // 8 successors
		nil, // No delegate
		160, // 160bit hash function
	}
}

// Computes the forward distance from a to b modulus a ring size
func distance(a, b []byte, bits int) *big.Int {
	// Get the ring size
	var ring big.Int
	ring.Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)

	// Convert to int
	var a_int, b_int big.Int
	(&a_int).SetBytes(a)
	(&b_int).SetBytes(b)

	// Compute the distances
	var dist big.Int
	(&dist).Sub(&b_int, &a_int)

	// Distance modulus ring size
	(&dist).Mod(&dist, &ring)
	return &dist
}

func TestDistance(t *testing.T) {
	a := []byte{63}
	b := []byte{3}
	d := distance(a, b, 6) // Ring size of 64

	fmt.Println("- testing distance function")
	fmt.Printf("expect distance 4 => result is %v \n", d)
	if d.Cmp(big.NewInt(4)) != 0 {
		t.Fatalf("expect distance 4! %v", d)
	}

	a = []byte{0}
	b = []byte{65}
	d = distance(a, b, 7) // Ring size of 128

	fmt.Printf("expect distance 65 => result is %v \n", d)

	if d.Cmp(big.NewInt(65)) != 0 {
		t.Fatalf("expect distance 65! %v", d)
	}

	a = []byte{1}
	b = []byte{255}
	d = distance(a, b, 8) // Ring size of 256

	fmt.Printf("expect distance 254 => result is %v \n", d)
	if d.Cmp(big.NewInt(254)) != 0 {
		t.Fatalf("expect distance 254! %v", d)
	}
}

/*func main() {
	fmt.Println("testing")
	conf := DefaultConfig("localhost")
	//v := &Vnode{Id: []byte{59}, "localhost"}
	v := &Vnode{}

	idx := 10
	v.genId(uint16(idx), conf)
	fmt.Println(v)
} */
