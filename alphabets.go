package key

// Ready-made alphabets for common cases. They are ordinary strings: pass any
// of them to NewDynamic or NewFixed, or build your own. Ordering defines the
// digit values, so two alphabets with the same characters in a different order
// produce different keys.
const (
	// Base16 is the lowercase hexadecimal alphabet (0-9a-f).
	Base16 = "0123456789abcdef"

	// Base32 is the RFC 4648 base32 alphabet (A-Z2-7), without padding.
	Base32 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

	// Base36 is digits followed by lowercase letters (0-9a-z).
	Base36 = "0123456789abcdefghijklmnopqrstuvwxyz"

	// Base62 is digits, uppercase, then lowercase letters (0-9A-Za-z).
	Base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	// Crockford is Crockford's base32 alphabet: digits and uppercase letters
	// with I, L, O and U removed to avoid visual ambiguity.
	Crockford = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

	// Unambiguous drops characters that are easy to confuse when typed by hand
	// (0/O, 1/I/l). Handy for tickets, coupons and one-time codes.
	Unambiguous = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"
)
