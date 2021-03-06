package main
import (
	_"fmt"
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

//from the cauchy matrix mat, select only rows from row_indexes
func create_cauchy_submatrix(mat [][]byte, row_indexes []int) [][]byte {
	n := len(mat[0])
	submat := make([][]byte, n) //create matrix for cauchy
	for i := range submat{
		submat[i] = make([]byte, n)
	}

	for i := range submat { //populate it with rows from whole cauchy matrix
		submat[i] = mat[row_indexes[i]][:]
	}

	return submat
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

//create an inverse of the cauchy submatrix corresponding to row indexes in row_indexes.
func create_inverse(mat [][]byte, row_indexes []int) [][]byte {
	cauchy := create_cauchy_submatrix(mat, row_indexes)
	get_LU(cauchy)
	return invert_LU(cauchy)
}

func decode_word(inv [][]byte, enc []byte) []byte{
	dim := len(inv[0])

	//calculate W := (L^-1)[enc]
	w := make([]byte, dim)
	for r:=0; r<dim; r++ { //for every row in inv
		for j:=0; j<=r; j++ {
			if r == j { //diagonal values were overwritten in LU, but pretend they're still 1
				w[r] = add(w[r], enc[j])
			} else {
				w[r] = add(w[r], mul(inv[r][j], enc[j]))
			}
		}
	}

	data_word := make([]byte, dim)
	for r:=dim-1; r>=0; r-- {
		for j:=dim-1; j>=r; j-- {
			data_word[r] = add(data_word[r], mul(inv[r][j], w[j]))
		}
	}
	return data_word
}

