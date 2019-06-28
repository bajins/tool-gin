import sys
import random


def RandomBytes(n: int, no_zero_byte: bool = False):
    return bytes((random.randint(1 if no_zero_byte else 0, 255) for i in range(n)))


# from https://en.wikibooks.org/wiki/Algorithm_Implementation/Mathematics/Extended_Euclidean_algorithm
# return (g, x, y) where g = gcd(a, b) and g == a * x + b * y
def xgcd(b, a):
    x0, x1, y0, y1 = 1, 0, 0, 1
    while a != 0:
        q, b, a = b // a, a, b % a
        x0, x1 = x1, x0 - q * x1
        y0, y1 = y1, y0 - q * y1
    return b, x0, y0


def PKCS1_Padding(b: bytes, is_private_key_op: bool, sizeof_n: int):
    if len(b) > sizeof_n - 11:
        raise OverflowError('Message is too long.')

    ret = b'\x00\x01' if is_private_key_op else b'\x00\x02'
    ret += b'\xff' * (sizeof_n - 3 - len(b)) if is_private_key_op else RandomBytes(sizeof_n - 3 - len(b), True)
    ret += b'\x00'
    ret += b
    return ret


def PKCS1_Unpadding(b: bytes, sizeof_n: int):
    if len(b) != sizeof_n:
        raise ValueError('Message\'s length is not correct')

    if b.startswith(b'\x00\x01'):
        is_private_key_op = True
    elif b.startswith(b'\x00\x02'):
        is_private_key_op = False
    else:
        # I know it is also valid if b starts with b'\x00\x00',
        # but now I do not care about this situation.
        raise ValueError('It is not a PKCS1-padded message.')

    msg_start_ptr = 3
    while msg_start_ptr < len(b):
        if is_private_key_op and b[msg_start_ptr] == 0:
            break
        if is_private_key_op and b[msg_start_ptr] != 0xff:
            raise ValueError('It is not a PKCS1-padded message.')
        if not is_private_key_op and b[msg_start_ptr] == 0:
            break
        msg_start_ptr += 1
    msg_start_ptr += 1

    msg = b[msg_start_ptr:]
    if len(msg) > sizeof_n - 11:
        raise OverflowError('Message is too long.')

    return msg


def RSA_Encrypt(m: bytes, e: int, n: int):
    m = int.from_bytes(m, 'big')
    if m >= n:
        raise ValueError('Message is too big.')

    c = pow(m, e, n)

    return c.to_bytes((n.bit_length() + 7) // 8, 'big')


def RSA_Decrypt(c: bytes, d: int, n: int):
    c = int.from_bytes(c, 'big')
    if c >= n:
        raise ValueError('Ciphertext is too big.')

    m = pow(c, d, n)

    return m.to_bytes((n.bit_length() + 7) // 8, 'big')


p = 0x3862bf704e31d0962c0f27303efe8f5ba8d1edc08530351884522d3c1ddf289f
q = 0x3cd9629192d2a4b0645103b892b32901801770269e10b00e562ec34d817bd0fd
n = p * q
phi = (p - 1) * (q - 1)
e = 65537
d = xgcd(e, phi)[1]
while d < 0:
    d += phi


def GenLicenseCode(name: str, license_count: int):
    if license_count > 500 or license_count < 1:
        raise ValueError('Invalid license count.')

    info = '%s\n%d\n' % (name, license_count)
    msg = info.encode() + RandomBytes(4)
    padded_msg = PKCS1_Padding(msg, True, (n.bit_length() + 7) // 8)
    enc_msg = RSA_Encrypt(padded_msg, d, n)
    return enc_msg.hex()


if __name__ == '__main__':
    msg = GenLicenseCode("woytu", int(sys.argv[1]))
    print(msg)
