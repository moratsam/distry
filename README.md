# Distry

Implementation of some algos, math concepts


## Note: This is the current layout of the repo:

*	**cluster** contains my work on writing distributed system from the ground up, which I will eventually use to test other projects within this repo.
I am moving away from *libp2p* because deploying *libp2p* nodes to kubernetes leads to double-NAT-ing, which fucks with *libp2p*'s service discovery.

* **erasure_codes** is a stand-alone project, contains everything dealing with the implementation of Reed-Solomon erasure codes.

* **ssecret_sharing** is a stand-alone project, contains everything dealing with the implementation of Shamir's secret sharing.

* **everything else** is a stand-alone project, containing the implementation of nodes which depend on *libp2p*. This can be used to run Bracha's reliable broadcast, which is implemented in **rbc0*.


### HOWTO

### Note: this section deals with the **everything else** part of the repo


Execute the bootstrap bash script. It will use the distry-bootstrap.privkey to generate its identity. Then you can use the bash commands in the **bash** dir to spin up additional nodes and execute some services. You may have to change the IP in the bash scripts to match the IP of your bootstrap node.

#### repo structure overview

##### api
	
Implements the connection between grpc and the code.

##### bash

Utility commands to start a node etc.

##### cmd

The entry point of the program (the main function).

##### k8s

Some yamls for deployment to kubernetes. Not working yet because of double-NAT incompatibility with libp2p peer-discovery.

##### messages

Defines the structs used as messages in various protocols, as well as the logic for (un)marshalling (from)to protobuf structs.

##### node

Base logic of a node. Provides things like bootstrapping, identity creation and connection establishment.

##### omni

Implements the base medium through which nodes exchange messages (libp2p's pubsub).

##### proto

The protobuf definitions of services and messages.

##### proto\_gen

The auto-generated proto code.

##### reliable broadcast

Code for Bracha's reliable broadcast.

## rbc0

bracha's reliable broadcast
DOI: 10.1016/0890-5401(87)90054-x

We consider the following model of a distributed system. The system consists of n processes that communicate by sending messages through a message system. We assume a reliable message system in which no messages are lost or generated. Each process can directly send messages to any other process, and can identify the sender of every message it receives. Up to t of the processes are faulty and may deviate from the protocol. A protocol is called t-resilient if it satisfies the agreement and validity requirements in the presence of up to t faulty processes.

A protocol is a reliable broadcast protocol (*rbc*) if:
1. If process *p* is correct, then all correct processes agree on the value of its messsage;
2. If *p* is faulty than, either all correct processes agree on the same value or none of them accepts any value from *p*.

#### protocol 

The following is a *rbc* with 0 <= t < n/3 byzantine faulty processes.

***Broadcast(v)*** 
- **step 0.** (By process p )
	- Send *(initial,v)* to all the processes
	
- **step 1.** Wait until the receipt of
		one *(initial,v)* message
	 	or (n-t) *(echo,v )* messages
		or (t+1) *(ready,v)* messages
		for some v
	- Send *(echo,v)* to all the processes.

- **step 2.** Wait until the receipt of
		(n-t) *(echo,v)* messages
		or t+1 *(ready,v)* messages
		(including messages received in step 1)
		for some v
	- Send *(ready,v)* to all the processes.

- **step 3.** Wait until the receipt of,
		2t+1 *(ready,v)* messages
		(including messages received in step 1 or step 2) for some v.
	- Accept v.


* Lemma 1: If two correct processes *s* and *t* send *(ready, v)* and *(ready, u)* messages, respectively, then *u*=*v*.

PROOF: Let *q* be the first process that sends *(ready, v)* and *r* the first that sends *(ready, u)*. This means *q* must have received >= (n-t) *(echo, v)* messages and *r* must have received >= (n-t) *(echo, u)* messages. Intersection between two (n-t) subsets must includeat least (n-t)-t >= (t+1) elements which means at least one non-faulty process must have sent bot a *(ready, v)* and a *(ready, u)* message. But correct processes can send only one message of each type during a broadcast, hence a contradiction.


* Lemma 2: If two correct processes *q* and *r* accept the values *v* and *u*, respectively, then *u* = *v*.

PROOF: If *q* accepts *v* it must have received >= (2t+1) *(ready, v)* messages, at least (t+1) of which must have come from correct processes. Analogously for *r* and *u*. Hence, by lemma 1, *u* = *v*.


* Lemma 3: If a correct process *q* accepts a value *v* then every other correct process will eventually accept *v*

PROOF: For *q* to accept *v* it must gave received >= (2t+1) *(ready, v)* messages, of which at least (t+1) must have come from correct processes. Which means every process will eventually receive (t+1) *(ready, v)* messages, which means every correct process will eventually issue a *(ready, v)* message. Which means every correct process will eventually receive at least n-t >= 2t+1 *(ready, v)* messages and will thus accept *v*.


* Lemma 4: If a correct process *p* broadcasts *v* then all correct processes accept *v*.

PROOF: Every correct process *q* receives an *(init, v)* message and sends a *(echo, v)* message. Thus every correct process *q* will receive >= n-t *(echo, v)* messages and will send a *(ready, v)* message. Every correct process will receive >= n-t *(ready, v)* messages and will accept *v*. 

## erasure codes
*Polynomial Codes over Certain Finite Fields*
DOI: 10.1137/0108018

*Optimizing Cauchy Reed-Solomon Codes for Fault-Tolerant Storage Applications* 
DOI: 10.1.1.140.2267

#### Field
set of elements with (+, \*)
with (+, \*) identities
with (+, \*) inverses
division by id(+) not defined.


#### Finite field Zp:
* Lemma 1: Rows in permutation table except row 0 are permutations of [p-1].

PROOF: Suppose not. Suppose x*a = x*b. Then x(a-b) = 0 is divisor of zero. //

#### Galois field
	GF(p^m) are polynomials of degree m-1 over Zp. For example, ax^m-1 + bx^m-2 +..+ f where {a,..f} in [p-1]. 
	Addition and multiplication of the coefficients (but not the polynomials) are defined by Zp.
		addition table for Z2 (XOR)
			+	0	1
			0	0	1
			1	1	0

		multiplication table for Z2 (AND)
			*	0	1
			0	0	0
			1	0	1

	Problem seems to arise: multiplication on polynomials is not closed.

	A **prime** for GF(p^m) is a degree m polynomial that is irreducible over p . This simply means that it cannot be factored. For example, x^3 + 1 is not irreducible over 2 because it can be factored as (x^2 + x + 1)(x + 1).
	If an irreducible polynomial g(x) can be found, then polynomial multiplication can be defined as standard polynomial multiplication modulo g(x).
```
Example for GF(2^3), g(x) = x^3 + x + 1

Dec	Bin	Poly

0		000	0
1		001	1
2		010	x
3		011	x + 1
4		100	x^2
5		101	x^2 + 1
6		110	x^2 + x
7		111	x^2 + x + 1

5*6 != 30 % 8 = 6
5*6 = (x^2 + 1)(x^2 + x) % x^3 + x + 1 = x + 1 = 3
```

#### Galois field arithmetic

GF(2^k) addition or subtraction is xor.
To multiply *a* with *b*, imagine the binary written form as a polynomial of some *x* over {0,1}. Wherever there is a '1' in *a* it means add to the final result that power of *x* multiplied by *b*. Which of course translates to just right shift b by that power. This is done for each '1' in *a*. And how are these partial results then added together? Still thinking of the polynomial representation, it becomes obvious that the simply need to be summed up which is just XOR. Thus, multiplication can be easily implemented with a series of bit shifts and XORs.

The outcome of this operation must by divided by the prime polynomial to ensure that the end result remains in GF(2^k). Now, thinking again in terms of polynomials, division is just subtraction of the divisor at the appropriate powers. And subtraction is also just XOR. The process stops once the the remainder is under 2^k, because for every *e* in GF(2^k): e divided by the prime is *e*.

example: a=33, b=191, prime=0x11d
```
  00100001 #a
* 10111111 #b
  =================
   _____10111111 #this is the rightmost '1' in a; the free coefficient in the polynomial so just *b* multiplied by 1
^ 10111111_____ #this is the second '1'. Here x is raised to the power of 5 so just shift b 5 times. 
  1011101011111 #normal multiplication is finished, result exceeds 2^8 -1


	
  # I try to align the divisor with the first 1 from the dividend, to ease understanding
  1011101011111 #it needs to be divided by the prime to arrive back in GF(2^8)
/ 100011101____ #this is 0x11d
  =================
  0011010001111 #still > 2^8 -1, so repeat
^   100011101__ #I al
  0001011111011 #repeat
^    100011101_
  0000011000001 # = 193 < 2^8 -1, end
```

A generator *g* of a field is an element of the field such that every other element of the field can be expressed as a series of iterative multiplications of *g*. In this way, *g* is said to generate the field.
To optimise multiplication, one can keep in memory the log and exp tables of a generator. Any multiplication in the field can then be performed by two lookups into the log table and 1 lookup into the exp table:

	a*b = g^(logg(a*b)) = g^(logg(a) + logg(b))
	
#### Reed-Solomon
Suppose we have a file of data of size N and suppose we want to create from it (n+k) shards of size roughly N/n, such that possessing any n-subset of the shards allows one to reconstruct the original data. 

note: all operations are performed on Galois fields (in my case GF(2^8) since I'm operating on bytes)

ENCODING: 
i) Create a cool matrix mat of dimensions (n+k)xn.
ii) Divide the data into words of size n (n-words) and stack them into a matrix [data] with dimensions nx(N/n)
iii) Define mat * [data] = [enc] //dimensions of [enc] are obviously [n+k]x[N/n]
Each row in [enc] can be thought of as a shard.
The index of the row should be put into the shard, as it's needed for decoding.


DECODING:
This is the magic idea:
	mat * [data] = [enc] ==> 
	==> mat^-1 * mat * [data] = mat^-1 * [enc] ==>
	==> [data] = mat^-1 * [enc]


Say one has n shards. 
i) stack them together to create a submatrix of [enc] called [subenc]
ii) Create mat
iii) Remove all rows not pertaining to one's shards // now one is left with a nxn submat
iv) Calculate submat^-1, the inverse of submat
v) Reconstruct [data] by multiplying submat^-1 * [subenc]


The main difficulty with these scheme is that mat must have the property that every possible submat must be invertible. I used a standard cauchy matrix for this purpose. Authors recommend appending an identity matrix to the top, to cleanly separate the encoded data into data shards and parity shards. I disregarded this and used a complete cauchy matrix so every shard is encoded.

I implemented the matrix inversion using LU decomposition, of course with the twist that matrix values are polynomials over GF(2^8) and all operations also take place in that field.


![ec1](https://user-images.githubusercontent.com/43090095/138609366-a6258490-2764-4a07-8c83-9b28ee44b800.jpg)

![ec2](https://user-images.githubusercontent.com/43090095/138609362-78657857-64a6-4241-ac90-8cefb1d2dd4f.jpg)


Some notes:

* Any sumbatrix of my original cauchy matrix is a viable cauchy matrix.
* Picture 1 rewrites a cauchy matrix as a product using vandermonde matrixes.
* Picture 2 proves that a cauchy matrix is invertible by providing a formula for the determinant of a (square) cauchy matrix.
* The derivation defines a cauchy matrix as subtracting different indexed items, whereas I implemented it as additions. The derivation still applies because under GF subtraction and addition are equivalent.
* Since x and y terms are pairwise-disjunct by definition, the derived determinant will always be different from zero, which means it is indeed invertible.


### Code

* io.go implements parallel reading and writing to files

* gf\_arithmetic.go implements add, sub, mul, div over GF(2^8)

* reed\_solomon.go contains the heart of the project. It implements cauchy matrix creation, LU decomposition and matrix inversion and the decoding of encoded data.

* manager.go represents the out-facing side of the project. It provides a Manager through which one can invoke the Encode and Decode functions to perform reed-solomon encoding / encoding.


###### Encode
A cauchy matrix is constructed, when the Manager is initated and before Encode() is even called.

Takes as argument the path to a file which seeks to be encoded.
Returns the filepaths to the encoded shards.

A routine reads the file with input data and passes it on to the encoder.
The encoder spawns (n+k) sub-encoding routines - one for each row of the cauchy matrix.
Taking words of size n (n-words) from the incoming data, each routine calculates the dot product between the n-word and its row of the cauchy matrix.
Then each sub-encoder sends the data to a file writer routine, which writes the encoded data to a file - for a total of (n+k) shards.

**See the code for more detailed comments**

###### Decode

Takes 2 arguments: the filepaths to the shards and a name for the output file (which will be created).

First, all the shards are opened and their first byte read to find the indexes of the cauchy matrix rows that were used to create them.
These indexes are then used to create the appropriate cauchy sub-matrix.
An inverse of the sub-matrix is then calculated using LU decomposition.
Then, data is passed from the shards, one word at a time and is decoded using the inverted sub-matrix.
The decoded data is passed to the writer routine.

**See the code for more detailed comments**


## Shamir's secret sharing

One has a secret that one wishes to share between *n* people.
However, none of these *n* people can find the secret by themselves.
Instead, any *k* of these people must come together and combine their secrets to obtain the final secret.

Super easy to understand: Construct a random polynomial *f(x)* of degree *(k-1)* and insert the secret as its constant term (the coefficient next to x^0). As per Lagrange's theorem, given any *k* points **(xi, f(xi))** one can reconstruct the polynomial f(x).
So simply evaluate the polynomial at some *n* pre-determined points (e.g. from 1 to n) and give each of the *n* people their pair **(xi, f(xi)**.


Super easy to implement, since I reused the Galois Field arithmetic code originally written for *erasure codes*.
The secret (constant term) was found using the standard Lagrange polynomial interpolation. Since for Shamir's secret sharing scheme one is only interested in the constant term, the creation of the Lagrange base polynomials can be optimised considerably (Basically, only the constant term of the base polynomials needs to be calculated).
