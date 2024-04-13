// Welcome to the gnark playground!
package main

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/mimc"
)

// gnark is a zk-SNARK library written in Go. Circuits are regular structs.
// The inputs must be of type frontend.Variable and make up the witness.
// The witness has a
//   - secret part --> known to the prover only
//   - public part --> known to the prover and the verifier
type Circuit struct {
	Secret frontend.Variable // pre-image of the hash secret known to the prover only
	Hash   frontend.Variable `gnark:",public"` // hash of the secret known to all
}

// Define declares the circuit logic. The compiler then produces a list of constraints
// which must be satisfied (valid witness) in order to create a valid zk-SNARK
// This circuit proves knowledge of a pre-image such that hash(secret) == hash
func (circuit *Circuit) Define(api frontend.API) error {
	// hash function
	mimc, _ := mimc.NewMiMC(api)

	// hash the secret
	mimc.Write(circuit.Secret)

	// ensure hashes match
	api.AssertIsEqual(circuit.Hash, mimc.Sum())

	return nil
}

func main() {
	// compiles our circuit into a R1CS
	var circuit Circuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(ccs)

	// witness definition
	assignment := Circuit{Secret: "0xdeadf00d", Hash: "1037254799353855871006189384309576393135431139055333626960622147300727796413"}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()
	fmt.Println(assignment)
	fmt.Println(witness)
	fmt.Println(publicWitness)

	// groth16: Prove & Verify
	proof, _ := groth16.Prove(ccs, pk, witness)
	groth16.Verify(proof, vk, publicWitness)

	// ccs needs to be constructed in the protocol
	//
}
