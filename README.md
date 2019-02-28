# OpenToken_Decoder
OpenToken (OTK), a format for the lightweight, secure, cross-application exchange of key-value pairs between applications that use HTTP (see [RFC2616]) as the transport protocol.  The format is designed primarily for use as an HTTP cookie (see [RFC2965]) or query parameter, but can also be used in other scenarios that require a compact, application-neutral token.

The OpenToken technology is not designed to encapsulate formal identity assertions (for which see [SAML]) or authentication credentials (for which see [SASL]).  Instead, OpenToken is designed to encapsulate basic name-value pairs for exchange between applications that use HTTP as the transport protocol.

This project is used to decode an OpenToken following the rules in Section 4 of [RFC4648], RFC1950, RFC1951, initilaising HMAC using the SHA-1 algorithm

The entire process basically involves as described at https://tools.ietf.org/html/draft-smith-opentoken-02
1.   Replace the "*" padding characters with standard Base 64 "=" characters.
2.   Base 64 decode the OTK, following the rules in Section 4 of
        [RFC4648] and ensuring that the padding bits are set to zero.
3.   Validate the OTK header literal and version.
4.   Extract the Key Info (if present) and select a key for
        decryption.
5.   Decrypt the payload cipher-text using the selected cipher suite.
6.   Decompress the decrypted payload, in accordance with [RFC1950]
        and [RFC1951].
7.   Initialize an [HMAC] using the SHA-1 algorithm specified in
        [SHA] and the following data (order is significant):
        1.  OTK version
        2.  Cipher suite value
        3.  IV value (if present)
        4.  Key Info value (if present)
        5.  Payload length (2 bytes, network order)
8.   Update the HMAC from the previous step with the clear-text
        payload (after decompressing).
9.   Compare the HMAC from step 8 with the HMAC received in the OTK.
        If they do not match, halt processing.
10.  Process the payload.
