// Package key provides tools for generating and decoding unique keys
// based on a given alphabet.
//
// The main structure of the package is Locksmith, which generates keys
// depending on the iteration number. This allows you to receive unique keys
// from a given alphabet without having to store all generated keys.
//
// Package provides flexibility regarding key size. The size can be set at
// the Locksmith creation stage or used dynamic size. When using dynamic
// key size, its length will change accordingly of the iteration number,
// but will not exceed the maximum size set.
//
// To guarantee the uniqueness of the keys, the alphabet from which the keys
// are generated must not contain repeated characters.
package key
