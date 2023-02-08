# httpmux

I realized the multiplexing to the right service is completely unrelated. Given
that we know that the only way to write correct programs is to keep them small I
decided to move the HTTP multiplexing and the SEAL certificate handling outside
of ergoproxy. Until I look into how to handle the use cases described above I'll
use nginx.

- https://pkg.go.dev/golang.org/x/crypto/acme
