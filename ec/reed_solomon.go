package main
import (
	"fmt"
)

//------------------------------------

//create cauchy matrix of dimensions (n+k)xn
//every n rows suffice to reconstruct the data
func create_cauchy(k, n byte) [][]byte{
	mat := make([][]byte, n+k)
	for i := range mat {
		mat[i] = make([]byte, n)
	}

	var i, j byte
	for i=0; i<n+k; i++ {
		for j=n+k; j<2*n+k; j++ {
			mat[i][j-n-k] = div(1, add(i, j))
		}
	}
	return mat
}

// enc encoded, (n+k)x1,  d: data, (n)x1
//e = (m.mat)*d
func encode(data []byte, mat [][]byte) []byte {
	k, n := byte(len(mat)-len(mat[0])), byte(len( mat[0] ))
	fmt.Println()
	
	enc := make([]byte, n+k)
	var i, j byte
	for i=0; i<n+k; i++ {
		for j=0; j<n; j++ {
			enc[i] = add(enc[i], mul(mat[i][j], data[j]))
		}
	}
	return enc
}


func get_LU(mat [][]byte) {
	dim := byte(len( mat[0] ))

	var i, row_ix, col_ix byte
	for i=0; i<dim; i++{
		if mat[i][i] == 0{
			continue
		}
		for row_ix=i+1; row_ix<dim; row_ix++{
			//derive factor to destroy first elemnt
			mat[row_ix][i] = div(mat[row_ix][i], mat[i][i])
			//subtract (row i's element * factor) from every other element in row
			for col_ix=i+1; col_ix<dim; col_ix++{
				mat[row_ix][col_ix] = sub(mat[row_ix][col_ix], mul(mat[i][col_ix],mat[row_ix][i]))
			}
		}
	}
}

func invert_LU(mat [][]byte) [][]byte {
	dim := len( mat[0] )

	side := make([][]byte, dim) //create side identity matrix
	for i := range side {
		side[i] = make([]byte, dim)
		side[i][i] = 1
	}

	//invert U by adding an identity to its side. When U becomes identity, side is inverted U.
	//no operations on U actually need to be performed, just their effects on the side
	//matrix are being recorded.
	var i, j, k int
	for i=dim-1; i>=0; i-- { //for every row
		for j=dim-1; j>i; j-- { //for every column
			for k=dim-1; k>=j; k-- { //subtract row to get a 0 in U, reflect this change in side
				side[i][k] = sub(side[i][k], mul(mat[i][j], side[j][k]))
			}
		}
		if mat[i][i] == 0{
			continue
		} else {
			//divide mat[i][i] by itself to get a 1, reflect this change in whole line of side
			for j=dim-1; j>=0; j-- {
				side[i][j] = div(side[i][j], mat[i][i])
			}
		}
	}

	//get inverse of L
	for i=0; i<dim; i++ {
		for j=0; j<i; j++ {
			for k=0; k<=j; k++ {
				//since an in-place algo was used for LU decomposition,
				//diagonal values of LU were overwritten by U,
				//whereas L expects them to be equal to 1
				//in this case, no mul should be performed (to simulate multiplying by 1)
				if j == k { 
					side[i][k] = sub(side[i][k], mat[i][j])
				} else {
					side[i][k] = sub(side[i][k], mul(mat[i][j], side[j][k]))
				}
			}
		}
	}

	//inverse matrix is now the side matrix! because m.inv kinda became identity matrix
	//kinda, because no changes to m.inv were actually recorded
	return side
}


//(mat)[data] = [enc]
//(mat^-1)(mat)[data] = (mat^-1)[enc] ==> (mat^-1)[enc] = [data]
//mat = LU
//mat^-1 = (U^-1)(L^-1)  
//(U^-1)(L^-1)[enc] = [data]
func solve_from_inverse(inv, tmp_enc [][]byte) []byte {
	enc := []byte{tmp_enc[0][1], tmp_enc[1][1], tmp_enc[2][1]}

	dim := len( inv[0] )
	var i, j int

	//calculate W := (L^-1)[enc]
	w := make([]byte, dim)
	for i=0; i<dim; i++ {
		for j=0; j<=i; j++ {
			if i == j { //diagonal values were overwritten, but pretend they're still 1
				w[i] = add(w[i], enc[j])
			} else {
				w[i] = add(w[i], mul(inv[i][j], enc[j]))
			}
		}
	}

	//calculate [data] = (U^-1)W
	data := make([]byte, dim)
	for i=dim-1; i>=0; i-- {
		for j=dim-1; j>=i; j-- {
			data[i] = add(data[i], mul(inv[i][j], w[j]))
		}
	}
	return data
}