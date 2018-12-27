/*
 * Copyright (c) 2018 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 *
 * This application uses Open Source components. You can find the
 * source code of their open source projects along with license
 * information below. We acknowledge and are grateful to these
 * developers for their contributions to open source.
 *
 * Project: CockroachDB https://github.com/cockroachdb/cockroach
 * Copyright 2018 The Cockroach Authors.
 * License (Apache License 2.0) https://github.com/cockroachdb/cockroach/blob/master/LICENSE
 *
 * Project: Vitess https://github.com/vitessio/vitess
 * Copyright 2018 Google Inc.
 * License (Apache License 2.0) https://github.com/vitessio/vitess/blob/master/LICENSE
 *
 * Project: Citus https://github.com/citusdata/citus
 * Copyright 2018 Citus Data, Inc.
 * License (GNU Affero General Public License v3.0) https://github.com/citusdata/citus/blob/master/LICENSE
 *
 * Project: pg_query_go https://github.com/lfittl/pg_query_go
 * Copyright 2018 Lukas Fittl
 * License (3-Clause BSD) https://github.com/lfittl/pg_query_go/blob/master/LICENSE
 *
 * Project: pgx https://github.com/jackc/pgx
 * Copyright 2018 Jack Christensen
 * License (MIT) https://github.com/jackc/pgx/blob/master/LICENSE
 *
 * Project: BadgerDB https://github.com/dgraph-io/badger
 * Copyright 2018 Dgraph Labs, Inc. and Contributors
 * License (MIT) https://github.com/dgraph-io/badger/blob/master/LICENSE
 *
 * Project: Sonyflake https://github.com/sony/sonyflake
 * Copyright 2018 Sony Corporation
 * License (MIT) https://github.com/sony/sonyflake/blob/master/LICENSE
 *
 * Project: Raft https://github.com/hashicorp/raft
 * Copyright 2018 HashiCorp
 * License (MPL-2.0) https://github.com/hashicorp/raft/blob/master/LICENSE
 *
 * Project: pq github.com/lib/pq
 * Copyright 2018  'pq' Contributors Portions Copyright (C) 2018 Blake Mizerany
 * License https://github.com/lib/pq/blob/master/LICENSE.md
 *
 * Project: go-linq github.com/ahmetb/go-linq
 * Copyright 2018 Ahmet Alp Balkan
 * License https://github.com/ahmetb/go-linq/blob/master/LICENSE
 */

package ipaddr

import (
    "bytes"
    "encoding/binary"
    "io"
    "math"
    "math/rand"
    "net"
    "strconv"
    "strings"

    "github.com/readystock/noah/db/sql/pgwire/pgerror"
    "github.com/readystock/noah/db/util/uint128"
)

var errResultOutOfRange = pgerror.NewError(pgerror.CodeNumericValueOutOfRangeError, "result out of range")

// Addr is the representation of the IP address. The Uint128 takes 16-bytes for
// both IPv4 and IPv6.
type Addr uint128.Uint128

// IPAddr stores an IP address's family, IP, and host mask. This was chosen over
// Go's "net" IP, as that struct doesn't work well for what we need to do.
//     - It discards information when parsing IPv4, forcing it to be IPv6, and then
//         assuming IPv4-mapped IPv6 addresses are purely IPv4 (only for printing).
//         This is solved by having a Family field.
//     - ParseIP and ParseCIDR are very strict, whereas postgres' INET and CIDR
//       have very relaxed constraints for parsing an IP.
//  - Doing int64 operations is much more efficient than byte slice operations.
type IPAddr struct {
    // Family denotes what type of IP the original IP was.
    Family IPFamily
    Mask   byte
    Addr   Addr
}

// IPFamily denotes which classification the IP belongs to.
type IPFamily byte

const (
    // IPv4family is for IPs in the IPv4 space.
    IPv4family IPFamily = iota
    // IPv6family is for IPs in the IPv6 space.
    IPv6family
)

// IPv4size 1 byte for family, 1 byte for mask, 4 for IP.
const IPv4size = net.IPv4len + 2

// IPv6size 1 byte for family, 1 byte for mask, 16 for IP.
const IPv6size = net.IPv6len + 2

// IPv4mappedIPv6prefix is the byte prefix for IPv4-mapped IPv6.
const IPv4mappedIPv6prefix = uint64(0xFFFF) << 32

// IPv4mask is used to select only the lower 32 bits of uint128.
// IPv4 addresses may not have the the upper 96 bits of Addr set to 0.
// IPv4 addresses mapped to IPv6 have prefix bits that should not change.
const IPv4mask = uint64(0xFFFFFFFF)

// IPv4max is used for overflows.
const IPv4max = IPv4mask

// ToBuffer appends the IPAddr encoding to a buffer and returns the final buffer.
func (ipAddr *IPAddr) ToBuffer(appendTo []byte) []byte {
    // Record the family as the first byte and the mask size as the second byte.
    appendTo = append(appendTo, byte(ipAddr.Family), ipAddr.Mask)
    if ipAddr.Family == IPv4family {
        appendTo = append(appendTo, 0, 0, 0, 0)
        binary.BigEndian.PutUint32(appendTo[len(appendTo)-4:], uint32(ipAddr.Addr.Lo))
    } else {
        appendTo = append(appendTo, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
        binary.BigEndian.PutUint64(appendTo[len(appendTo)-16:len(appendTo)-8], ipAddr.Addr.Hi)
        binary.BigEndian.PutUint64(appendTo[len(appendTo)-8:], ipAddr.Addr.Lo)
    }
    return appendTo
}

// FromBuffer populates an IPAddr with data from a byte slice, returning the
// remaining buffer or an error.
func (ipAddr *IPAddr) FromBuffer(data []byte) ([]byte, error) {
    ipAddr.Family = IPFamily(data[0])
    if ipAddr.Family != IPv4family && ipAddr.Family != IPv6family {
        return nil, pgerror.NewErrorf(pgerror.CodeInternalError, "IPAddr decoding error: bad family, got %d", ipAddr.Family)
    }
    ipAddr.Mask = data[1]

    if ipAddr.Family == IPv4family {
        ipAddr.Addr.Lo = uint64(binary.BigEndian.Uint32(data[2:])) | IPv4mappedIPv6prefix
        return data[IPv4size:], nil
    }
    ipAddr.Addr = Addr(uint128.FromBytes(data[2:]))
    return data[IPv6size:], nil
}

// String will convert the IP to the appropriate family formatted string
// representation. In order to retain postgres compatibility we ensure
// IPv4-mapped IPv6 stays in IPv6 format, unlike net.Addr.String().
func (ipAddr IPAddr) String() string {
    ip := net.IP(uint128.Uint128(ipAddr.Addr).GetBytes())
    isIPv4MappedIPv6 := ipAddr.Family == IPv6family && ip.Equal(ip.To4())
    var maskSize byte

    if ipAddr.Family == IPv4family {
        maskSize = 32
    } else {
        maskSize = 128
    }
    if ipAddr.Mask == maskSize && isIPv4MappedIPv6 {
        return "::ffff:" + ip.String()
    } else if ipAddr.Mask == maskSize {
        return ip.String()
    } else if isIPv4MappedIPv6 {
        // Due to an issue with IPv4-mapped IPv6 the mask is also reduced to an IPv4
        // mask during net.IPNet.String, so we need to add the mask manually.
        return "::ffff:" + ip.String() + "/" + strconv.Itoa(int(ipAddr.Mask))
    }
    return ip.String() + "/" + strconv.Itoa(int(ipAddr.Mask))
}

// Compare two IPAddrs. IPv4-mapped IPv6 addresses are not equal to their IPv4
// mapping. The order of order importance goes Family > Mask > IP-bytes.
func (ipAddr IPAddr) Compare(other *IPAddr) int {
    if ipAddr.Family < other.Family {
        return -1
    } else if ipAddr.Family > other.Family {
        return 1
    }

    if ipAddr.Mask < other.Mask {
        return -1
    } else if ipAddr.Mask > other.Mask {
        return 1
    }
    return ipAddr.Addr.Compare(other.Addr)
}

// Equal checks if the family, mask, and IP are equal.
func (ipAddr *IPAddr) Equal(other *IPAddr) bool {
    return ipAddr.Family == other.Family && ipAddr.Mask == other.Mask && ipAddr.Addr.Equal(other.Addr)
}

// getFamily checks what family the ip is in. If it doesn't appear to match
// either, getFamily returns -1.
func getFamily(addr string) IPFamily {
    // Get the family of the IP.
    for i := 0; i < len(addr); i++ {
        switch addr[i] {
        case '.':
            return IPv4family
        case ':':
            return IPv6family
        }
    }
    // Default to IPv4, as we need to handle '192/10'.
    return IPv4family
}

// ParseINet parses postgres style INET types. See TestIPAddrParseINet for
// examples.
func ParseINet(s string, dest *IPAddr) error {
    var maskSize byte
    i := strings.IndexByte(s, '/')
    family := getFamily(s)

    // If no mask suffix was provided, implicitly don't mask.
    if i < 0 {
        // Trims IPv4 suffix "." to match postgres compitibility.
        addr := s
        if family == IPv4family {
            addr = strings.TrimRight(addr, ".")
            maskSize = 32
        } else {
            maskSize = 128
        }
        ip := net.ParseIP(addr)
        if ip == nil {
            return pgerror.NewErrorf(pgerror.CodeInvalidTextRepresentationError, "could not parse %q as inet. invalid IP", s)
        }

        *dest = IPAddr{Family: family,
            Addr: Addr(uint128.FromBytes(ip)),
            Mask: maskSize,
        }
        return nil
    }

    addr, maskStr := s[:i], s[i+1:]
    // Trims IPv4 suffix "." to match postgres compitibility.
    if family == IPv4family {
        addr = strings.TrimRight(addr, ".")
    }
    maskOnes, err := strconv.Atoi(maskStr)
    if err != nil {
        return pgerror.NewErrorf(pgerror.CodeInvalidTextRepresentationError, "could not parse %q as inet. invalid mask", s)
    } else if maskOnes < 0 || (family == IPv4family && maskOnes > 32) || (family == IPv6family && maskOnes > 128) {
        return pgerror.NewErrorf(pgerror.CodeInvalidTextRepresentationError, "could not parse %q as inet. invalid mask", s)
    }

    if family == IPv4family {
        // If the mask is outside the defined octets, postgres will raise an error.
        octetCount := strings.Count(addr, ".") + 1
        if (octetCount+1)*8-1 < maskOnes {
            return pgerror.NewErrorf(pgerror.CodeInvalidTextRepresentationError, "could not parse %q as inet. mask is larger than provided octets", s)
        }

        // Append extra ".0" to ensure there are a total of 4 octets.
        var buffer bytes.Buffer
        buffer.WriteString(addr)
        for i := 0; i < 4-octetCount; i++ {
            buffer.WriteString(".0")
        }
        addr = buffer.String()

    }

    ip := net.ParseIP(addr)
    if ip == nil {
        return pgerror.NewErrorf(pgerror.CodeInvalidTextRepresentationError, "could not parse %q as inet. invalid IP", s)
    }

    *dest = IPAddr{Family: family,
        Addr: Addr(uint128.FromBytes(ip)),
        Mask: byte(maskOnes),
    }
    return nil
}

// RandIPAddr generates a random IPAddr. This includes random mask size and IP
// family.
func RandIPAddr(rng *rand.Rand) IPAddr {
    var ipAddr IPAddr
    if rng.Intn(2) > 0 {
        ipAddr.Family = IPv4family
        ipAddr.Mask = byte(rng.Intn(33))
        ipAddr.Addr = Addr(uint128.FromInts(0, uint64(rng.Uint32())|IPv4mappedIPv6prefix))
    } else {
        ipAddr.Family = IPv6family
        ipAddr.Mask = byte(rng.Intn(129))
        ipAddr.Addr = Addr(uint128.FromInts(rng.Uint64(), rng.Uint64()))
    }
    return ipAddr
}

// Hostmask returns the host masked IP. This is defined as the IP address bits
// that are not masked.
func (ipAddr *IPAddr) Hostmask() IPAddr {
    var newIPAddr IPAddr
    newIPAddr.Family = ipAddr.Family
    newIPAddr.Addr = ipAddr.Addr

    if ipAddr.Family == IPv4family {
        LoMask := ^uint32(0) >> ipAddr.Mask
        newIPAddr.Addr.Lo = uint64(LoMask) | IPv4mappedIPv6prefix
    } else if ipAddr.Mask <= 64 {
        newIPAddr.Addr.Hi = ^uint64(0) >> ipAddr.Mask
        newIPAddr.Addr.Lo = ^uint64(0)
    } else {
        newIPAddr.Addr.Hi = uint64(0)
        newIPAddr.Addr.Lo = ^uint64(0) >> (ipAddr.Mask - 64)
    }

    if newIPAddr.Family == IPv4family {
        newIPAddr.Mask = 32
    } else {
        newIPAddr.Mask = 128
    }

    return newIPAddr
}

// Netmask returns the network masked IP. This is defined as the IP address bits
// that are masked.
func (ipAddr *IPAddr) Netmask() IPAddr {
    var newIPAddr IPAddr
    newIPAddr.Family = ipAddr.Family
    newIPAddr.Addr = ipAddr.Addr

    if ipAddr.Family == IPv4family {
        LoMask := ^uint32(0) << (32 - ipAddr.Mask)
        newIPAddr.Addr.Lo = uint64(LoMask) | IPv4mappedIPv6prefix
    } else if ipAddr.Mask <= 64 {
        newIPAddr.Addr.Hi = ^uint64(0) << (64 - ipAddr.Mask)
        newIPAddr.Addr.Lo = uint64(0)
    } else {
        newIPAddr.Addr.Hi = ^uint64(0)
        newIPAddr.Addr.Lo = ^uint64(0) << (128 - ipAddr.Mask)
    }

    if newIPAddr.Family == IPv4family {
        newIPAddr.Mask = 32
    } else {
        newIPAddr.Mask = 128
    }

    return newIPAddr
}

// Broadcast returns a new IPAddr where the host mask of the IP address is a
// full mask, i.e. 0xFF bytes.
func (ipAddr *IPAddr) Broadcast() IPAddr {
    var newIPAddr IPAddr
    newIPAddr.Family = ipAddr.Family
    newIPAddr.Mask = ipAddr.Mask
    newIPAddr.Addr = ipAddr.Addr

    if newIPAddr.Family == IPv4family {
        LoMask := ^uint64(0) >> (32 + newIPAddr.Mask)
        newIPAddr.Addr.Lo = newIPAddr.Addr.Lo | LoMask
    } else if newIPAddr.Mask < 64 {
        LoMask := ^uint64(0)
        HiMask := ^uint64(0) >> newIPAddr.Mask
        newIPAddr.Addr.Lo = newIPAddr.Addr.Lo | LoMask
        newIPAddr.Addr.Hi = newIPAddr.Addr.Hi | HiMask
    } else {
        LoMask := ^uint64(0) >> (newIPAddr.Mask - 64)
        newIPAddr.Addr.Lo = newIPAddr.Addr.Lo | LoMask
    }

    return newIPAddr
}

// Complement returns a new IPAddr which is the bitwise complement of the
// original IP. Only the lower 32 bits are changed for IPv4.
func (ipAddr *IPAddr) Complement() IPAddr {
    var newIPAddr IPAddr

    newIPAddr.Family = ipAddr.Family
    newIPAddr.Mask = ipAddr.Mask

    var familyMask uint128.Uint128
    if newIPAddr.Family == IPv4family {
        familyMask = uint128.Uint128{Hi: 0, Lo: IPv4mask}
    } else {
        familyMask = uint128.Uint128{Hi: math.MaxUint64, Lo: math.MaxUint64}
    }
    newIPAddr.Addr = Addr(uint128.Uint128(ipAddr.Addr).Xor(familyMask))

    return newIPAddr
}

// And returns a new IPAddr which is the bitwise AND of two IPAddrs.
// Only the lower 32 bits are changed for IPv4.
func (ipAddr *IPAddr) And(other *IPAddr) (IPAddr, error) {
    var newIPAddr IPAddr
    if ipAddr.Family != other.Family {
        return newIPAddr, pgerror.NewErrorf(pgerror.CodeInvalidParameterValueError,
            "cannot AND inet values of different sizes")
    }
    newIPAddr.Family = ipAddr.Family

    if ipAddr.Mask > other.Mask {
        newIPAddr.Mask = ipAddr.Mask
    } else {
        newIPAddr.Mask = other.Mask
    }

    newIPAddr.Addr = ipAddr.Addr.and(other.Addr)

    return newIPAddr, nil
}

// Or returns a new IPAddr which is the bitwise OR of two IPAddrs.
// Only the lower 32 bits are changed for IPv4.
func (ipAddr *IPAddr) Or(other *IPAddr) (IPAddr, error) {
    var newIPAddr IPAddr
    if ipAddr.Family != other.Family {
        return newIPAddr, pgerror.NewErrorf(pgerror.CodeInvalidParameterValueError,
            "cannot OR inet values of different sizes")
    }
    newIPAddr.Family = ipAddr.Family

    if ipAddr.Mask > other.Mask {
        newIPAddr.Mask = ipAddr.Mask
    } else {
        newIPAddr.Mask = other.Mask
    }

    newIPAddr.Addr = ipAddr.Addr.or(other.Addr)

    return newIPAddr, nil
}

func (ipAddr *IPAddr) sum(o int64, neg bool) (IPAddr, error) {
    // neg carries information about whether to add or subtract other.
    // x - -y is the same as x + y, and x + -y is the same as x - y.
    var newIPAddr IPAddr
    newIPAddr.Family = ipAddr.Family
    newIPAddr.Mask = ipAddr.Mask

    var err error
    var o2 uint64
    if o < 0 {
        neg = !neg
        // PostgreSQL ver 10 seems to have an error with min int64.
        // It does not handle the overflow correctly.
        // maxip - minInt64 incorrectly returns a valid address.
        if o == math.MinInt64 {
            o2 = 1 << 63
        } else {
            o2 = uint64(-o)
        }
    } else {
        o2 = uint64(o)
    }

    if ipAddr.Family == IPv4family {
        if !neg {
            newIPAddr.Addr, err = ipAddr.Addr.ipv4Add(o2)
        } else {
            newIPAddr.Addr, err = ipAddr.Addr.ipv4Sub(o2)
        }
    } else {
        if !neg {
            newIPAddr.Addr, err = ipAddr.Addr.ipv6Add(o2)
        } else {
            newIPAddr.Addr, err = ipAddr.Addr.ipv6Sub(o2)
        }
    }

    return newIPAddr, err
}

// Add returns a new IPAddr that is incremented by an int64.
func (ipAddr *IPAddr) Add(o int64) (IPAddr, error) {
    return ipAddr.sum(o, false)
}

// Sub returns a new IPAddr that is decremented by an int64.
func (ipAddr *IPAddr) Sub(o int64) (IPAddr, error) {
    return ipAddr.sum(o, true)
}

// SubIPAddr returns the difference between two IPAddrs.
func (ipAddr *IPAddr) SubIPAddr(other *IPAddr) (int64, error) {
    var diff int64
    if ipAddr.Family != other.Family {
        return diff, pgerror.NewErrorf(pgerror.CodeInvalidParameterValueError,
            "cannot subtract inet addresses with different sizes")
    }

    if ipAddr.Family == IPv4family {
        return ipAddr.Addr.subIPAddrIPv4(other.Addr)
    }
    return ipAddr.Addr.subIPAddrIPv6(other.Addr)
}

func (ipAddr IPAddr) contains(other *IPAddr) bool {
    addrNetmask := ipAddr.Netmask()
    o, err := other.And(&addrNetmask)
    if err != nil {
        return false
    }
    t, _ := ipAddr.And(&addrNetmask)

    return t.Equal(&o)
}

// ContainsOrEquals determines if one ipAddr is in the same
// subnet as another or the addresses and subnets are equal.
func (ipAddr IPAddr) ContainsOrEquals(other *IPAddr) bool {
    return ipAddr.contains(other) && ipAddr.Mask <= other.Mask
}

// Contains determines if one ipAddr is in the same
// subnet as another.
func (ipAddr IPAddr) Contains(other *IPAddr) bool {
    return ipAddr.contains(other) && ipAddr.Mask < other.Mask
}

// ContainedByOrEquals determines if one ipAddr is in the same
// subnet as another or the addresses and subnets are equal.
func (ipAddr IPAddr) ContainedByOrEquals(other *IPAddr) bool {
    return other.contains(&ipAddr) && ipAddr.Mask >= other.Mask
}

// ContainedBy determines if one ipAddr is in the same
// subnet as another.
func (ipAddr IPAddr) ContainedBy(other *IPAddr) bool {
    return other.contains(&ipAddr) && ipAddr.Mask > other.Mask
}

// ContainsOrContainedBy determines if one ipAddr is in the same
// subnet as another or vice versa.
func (ipAddr IPAddr) ContainsOrContainedBy(other *IPAddr) bool {
    return ipAddr.contains(other) || other.contains(&ipAddr)
}

// WriteIPv4Bytes writes the 4-byte IPv4 representation. If the IP is IPv6 then
// the first 12-bytes are truncated.
func (ip Addr) WriteIPv4Bytes(writer io.Writer) error {
    return binary.Write(writer, binary.BigEndian, uint32(ip.Lo))
}

// WriteIPv6Bytes gets the 16-byte IPv6 representation.
func (ip Addr) WriteIPv6Bytes(writer io.Writer) error {
    err := binary.Write(writer, binary.BigEndian, ip.Hi)
    if err != nil {
        return pgerror.NewErrorf(pgerror.CodeInternalError, "%s", err)
    }
    return binary.Write(writer, binary.BigEndian, ip.Lo)
}

// Equal wraps the Uint128 equivilance.
func (ip Addr) Equal(o Addr) bool {
    return uint128.Uint128(ip).Equal(uint128.Uint128(o))
}

// Compare wraps the Uint128 equivilance.
func (ip Addr) Compare(o Addr) int {
    return uint128.Uint128(ip).Compare(uint128.Uint128(o))
}

// Sub wraps the Uint128 subtraction.
func (ip Addr) Sub(o uint64) Addr {
    return Addr(uint128.Uint128(ip).Sub(o))
}

// Add wraps the Uint128 addition.
func (ip Addr) Add(o uint64) Addr {
    return Addr(uint128.Uint128(ip).Add(o))
}

// String wraps net.IP.String().
func (ip Addr) String() string {
    return net.IP(uint128.Uint128(ip).GetBytes()).String()
}

// and wraps the Uint128 AND.
func (ip Addr) and(o Addr) Addr {
    return Addr(uint128.Uint128(ip).And(uint128.Uint128(o)))
}

// or wraps the Uint128 OR.
func (ip Addr) or(o Addr) Addr {
    return Addr(uint128.Uint128(ip).Or(uint128.Uint128(o)))
}

// ipv4Add adds uint64 to Uint128 and checks for IPv4 overflows.
func (ip Addr) ipv4Add(o uint64) (Addr, error) {
    if o > IPv4max || IPv4max-o < ip.Lo&IPv4mask {
        var newAddr Addr
        return newAddr, errResultOutOfRange
    }
    return ip.Add(o), nil
}

// ipv6Add adds uint64 to Uint128 and checks for IPv6 overflows.
func (ip Addr) ipv6Add(o uint64) (Addr, error) {
    newIP := ip.Add(o)
    if newIP.Compare(ip) < 0 {
        return newIP, errResultOutOfRange
    }
    return newIP, nil
}

// ipv4Add subtracts uint64 from Uint128 and checks for IPv4 overflows.
func (ip Addr) ipv4Sub(o uint64) (Addr, error) {
    if o > IPv4max || (ip.Lo&IPv4mask) < o {
        var newAddr Addr
        return newAddr, errResultOutOfRange
    }
    return ip.Sub(o), nil
}

// ipv6Sub subtracts uint64 from Uint128 and checks for IPv6 overflows.
func (ip Addr) ipv6Sub(o uint64) (Addr, error) {
    newIP := ip.Sub(o)
    if newIP.Compare(ip) > 0 {
        return newIP, errResultOutOfRange
    }
    return newIP, nil
}

// subIPAddrIPv4 adds two Uint128s to return int64 and checks for overflows.
func (ip Addr) subIPAddrIPv4(o Addr) (int64, error) {
    return int64(uint128.Uint128(ip).Lo&IPv4mask) - int64(uint128.Uint128(o).Lo&IPv4mask), nil
}

// subIPAddrIPv6 adds two Uint128s to return int64 and checks for overflows.
func (ip Addr) subIPAddrIPv6(o Addr) (int64, error) {
    var sign int64 = 1
    if ip.Compare(o) < 0 {
        ip, o = o, ip
        sign = -1
    }
    var lo, hi uint64
    lo = ip.Lo - o.Lo
    hi = ip.Hi - o.Hi
    if lo > ip.Lo {
        hi--
    }
    var diff int64
    if hi != 0 {
        return diff, errResultOutOfRange
    }
    if lo == uint64(math.MaxInt64)+1 && sign == -1 {
        return math.MinInt64, nil
    }
    if lo > uint64(math.MaxInt64) {
        return diff, errResultOutOfRange
    }
    diff = sign * int64(lo)
    return diff, nil
}
