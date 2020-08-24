package groth16

import (
	backend_bls377 "github.com/consensys/gnark/backend/bls377"
	backend_bls381 "github.com/consensys/gnark/backend/bls381"
	backend_bn256 "github.com/consensys/gnark/backend/bn256"
	"github.com/consensys/gnark/encoding"
	"github.com/consensys/gurvy"

	groth16_bls377 "github.com/consensys/gnark/backend/bls377/groth16"
	groth16_bls381 "github.com/consensys/gnark/backend/bls381/groth16"
	groth16_bn256 "github.com/consensys/gnark/backend/bn256/groth16"
	"github.com/consensys/gnark/backend/r1cs"
)

// Proof represents a Groth16 proof generated by groth16.Prove
// it's underlying implementation is curve specific (i.e bn256/groth16/Proof, ...)
type Proof interface{}

// ProvingKey represents a Groth16 ProvingKey
// it's underlying implementation is curve specific (i.e bn256/groth16/ProvingKey, ...)
type ProvingKey interface {
	IsDifferent(interface{}) bool
}

// VerifyingKey represents a Groth16 VerifyingKey
// it's underlying implementation is curve specific (i.e bn256/groth16/VerifyingKey, ...)
type VerifyingKey interface {
	IsDifferent(interface{}) bool
}

// Verify runs the groth16.Verify algorithm on provided proof with given solution
// it checks the underlying type of proof (i.e curve specific) to call the proper implementation
func Verify(proof Proof, vk VerifyingKey, solution map[string]interface{}) error {
	switch _proof := proof.(type) {
	case *groth16_bls377.Proof:
		return groth16_bls377.Verify(_proof, vk.(*groth16_bls377.VerifyingKey), solution)
	case *groth16_bls381.Proof:
		return groth16_bls381.Verify(_proof, vk.(*groth16_bls381.VerifyingKey), solution)
	case *groth16_bn256.Proof:
		return groth16_bn256.Verify(_proof, vk.(*groth16_bn256.VerifyingKey), solution)
	default:
		panic("unrecognized R1CS curve type")
	}
}

// Prove generate a groth16.Proof
// it checks the underlying type of the R1CS (curve specific) to call the proper implementation
func Prove(r1cs r1cs.R1CS, pk ProvingKey, solution map[string]interface{}) (Proof, error) {

	switch _r1cs := r1cs.(type) {
	case *backend_bls377.R1CS:
		return groth16_bls377.Prove(_r1cs, pk.(*groth16_bls377.ProvingKey), solution)
	case *backend_bls381.R1CS:
		return groth16_bls381.Prove(_r1cs, pk.(*groth16_bls381.ProvingKey), solution)
	case *backend_bn256.R1CS:
		return groth16_bn256.Prove(_r1cs, pk.(*groth16_bn256.ProvingKey), solution)
	default:
		panic("unrecognized R1CS curve type")
	}
}

// Setup runs groth16.Setup with provided R1CS
// it checks the underlying type of the R1CS (curve specific) to call the proper implementation
func Setup(r1cs r1cs.R1CS) (ProvingKey, VerifyingKey) {

	switch _r1cs := r1cs.(type) {
	case *backend_bls377.R1CS:
		var pk groth16_bls377.ProvingKey
		var vk groth16_bls377.VerifyingKey
		groth16_bls377.Setup(_r1cs, &pk, &vk)
		return &pk, &vk
	case *backend_bls381.R1CS:
		var pk groth16_bls381.ProvingKey
		var vk groth16_bls381.VerifyingKey
		groth16_bls381.Setup(_r1cs, &pk, &vk)
		return &pk, &vk
	case *backend_bn256.R1CS:
		var pk groth16_bn256.ProvingKey
		var vk groth16_bn256.VerifyingKey
		groth16_bn256.Setup(_r1cs, &pk, &vk)
		return &pk, &vk
	default:
		panic("unrecognized R1CS curve type")
	}
}

// DummySetup create a random ProvingKey with provided R1CS
// it doesn't return a VerifyingKey and is use for benchmarking or test purposes only.
func DummySetup(r1cs r1cs.R1CS) ProvingKey {
	switch _r1cs := r1cs.(type) {
	case *backend_bls377.R1CS:
		var pk groth16_bls377.ProvingKey
		groth16_bls377.DummySetup(_r1cs, &pk)
		return &pk
	case *backend_bls381.R1CS:
		var pk groth16_bls381.ProvingKey
		groth16_bls381.DummySetup(_r1cs, &pk)
		return &pk
	case *backend_bn256.R1CS:
		var pk groth16_bn256.ProvingKey
		groth16_bn256.DummySetup(_r1cs, &pk)
		return &pk
	default:
		panic("unrecognized R1CS curve type")
	}
}

// ReadProvingKey ...
// TODO likely temporary method, need a clean up pass on serialization things
func ReadProvingKey(path string) (ProvingKey, error) {
	curveID, err := encoding.PeekCurveID(path)
	if err != nil {
		return nil, err
	}
	var pk ProvingKey
	switch curveID {
	case gurvy.BN256:
		pk = &groth16_bn256.ProvingKey{}
	case gurvy.BLS377:
		pk = &groth16_bls377.ProvingKey{}
	case gurvy.BLS381:
		pk = &groth16_bls381.ProvingKey{}
	default:
		panic("not implemented")
	}

	if err := encoding.Read(path, pk, curveID); err != nil {
		return nil, err
	}
	return pk, err
}

// ReadVerifyingKey ...
// TODO likely temporary method, need a clean up pass on serialization things
func ReadVerifyingKey(path string) (VerifyingKey, error) {
	curveID, err := encoding.PeekCurveID(path)
	if err != nil {
		return nil, err
	}
	var vk VerifyingKey
	switch curveID {
	case gurvy.BN256:
		vk = &groth16_bn256.VerifyingKey{}
	case gurvy.BLS377:
		vk = &groth16_bls377.VerifyingKey{}
	case gurvy.BLS381:
		vk = &groth16_bls381.VerifyingKey{}
	default:
		panic("not implemented")
	}

	if err := encoding.Read(path, vk, curveID); err != nil {
		return nil, err
	}
	return vk, err
}

// ReadProof ...
// TODO likely temporary method, need a clean up pass on serialization things
func ReadProof(path string) (Proof, error) {
	curveID, err := encoding.PeekCurveID(path)
	if err != nil {
		return nil, err
	}
	var proof Proof
	switch curveID {
	case gurvy.BN256:
		proof = &groth16_bn256.Proof{}
	case gurvy.BLS377:
		proof = &groth16_bls377.Proof{}
	case gurvy.BLS381:
		proof = &groth16_bls381.Proof{}
	default:
		panic("not implemented")
	}

	if err := encoding.Read(path, proof, curveID); err != nil {
		return nil, err
	}
	return proof, err
}
