// Package key provides tools for generating and decoding unique, reversible
// string identifiers from numeric values using a custom alphabet. This package
// is particularly useful for creating URL-friendly identifiers, obfuscating
// sequential IDs, and generating human-readable unique keys.
//
// The package is built around the Locksmith type, which handles the conversion
// between numeric values and string keys. The conversion is bidirectional and
// deterministic - each numeric value maps to a unique string key, and each
// valid key maps back to its original numeric value.
//
// Common Use Cases:
//
//   - URL Shortening: Generate short, readable URLs from sequential IDs
//     Example:
//     ls, _ := key.New("abcdefghijklmnopqrstuvwxyz", 5)
//     shortURL, _ := ls.Marshal(1234567)
//
//   - ID Obfuscation: Hide sequential database IDs in public-facing identifiers
//     Example:
//     ls, _ := key.New("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 8)
//     publicID, _ := ls.Marshal(userID) // convert DB ID to public identifier
//     dbID, _ := ls.Unmarshal(publicID) // recover original DB ID when needed
//
//   - Ticket/Coupon Generation: Create unique, readable codes
//     Example:
//     ls, _ := key.New("23456789ABCDEFGHJKLMNPQRSTUVWXYZ", 6)
//     ticketCode, _ := ls.Marshal(ticketID)
//
//   - Resource Identifiers: Generate unique identifiers for resources
//     Example:
//     ls, _ := key.New("0123456789abcdef", 0) // dynamic length
//     resourceID, _ := ls.Marshal(timestamp)
//
// Key Features:
//
//   - Customizable alphabet for generated keys
//   - Support for both fixed and dynamic length keys
//   - Bidirectional conversion (numeric <-> string)
//   - No collision guarantee within the defined space
//   - Thread-safe operations
//   - No external dependencies
//
// The package ensures that generated keys are unique within the possible
// range determined by the alphabet length and key size. For fixed-size keys,
// the maximum possible value is alphabet_length^key_size. For dynamic-size
// keys, the maximum value is uint64.MaxValue.
//
// For optimal performance and security, consider:
//   - Choose an alphabet size appropriate for your use case.
//   - Avoid similar-looking characters if keys will be manually typed.
//   - Use fixed-size keys when possible for consistent length.
//   - Consider the trade-off between key length and alphabet size.
package key
