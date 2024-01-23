/*
* Rubrik is a document oriented database implementation with a rubrik's cube like datastructure.  
* Features are full CRUD operations, and full persistent to disk.  More coming..
* Copyright (C)  Alex Gaetano Padula
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"unsafe"
)

const cubeSize = 3

type KeyValueBlock struct {
	Key   string
	Value interface{}
}

type DocumentCube struct {
	File *os.File
	Size int
}

func NewDocumentCube(filePath string, size int) (*DocumentCube, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &DocumentCube{
		File: file,
		Size: size,
	}, nil
}

func (cube *DocumentCube) SetDocument(x, y, z int, document map[string]interface{}) error {
	position := cube.calculatePosition(x, y, z)
	enc := gob.NewEncoder(cube.File)

	if err := cube.File.Truncate(int64(position)); err != nil {
		return err
	}

	cube.File.Seek(int64(position), 0)

	if err := enc.Encode(document); err != nil {
		return err
	}

	return nil
}

func (cube *DocumentCube) GetDocument(x, y, z int) (map[string]interface{}, error) {
	position := cube.calculatePosition(x, y, z)
	dec := gob.NewDecoder(cube.File)

	info, err := cube.File.Stat()
	if err != nil {
		return nil, err
	}

	if info.Size() < int64(position)+int64(unsafe.Sizeof(KeyValueBlock{})) {
		return nil, fmt.Errorf("document not found")
	}

	cube.File.Seek(int64(position), 0)

	document := map[string]interface{}{}
	if err := dec.Decode(&document); err != nil {
		return nil, err
	}

	return document, nil
}

func (cube *DocumentCube) DeleteDocument(x, y, z int) error {
	position := cube.calculatePosition(x, y, z)
	dec := gob.NewDecoder(cube.File)
	enc := gob.NewEncoder(cube.File)

	info, err := cube.File.Stat()
	if err != nil {
		return err
	}

	if info.Size() < int64(position)+int64(unsafe.Sizeof(KeyValueBlock{})) {
		return fmt.Errorf("document not found")
	}

	// Read all documents after the deleted position
	cube.File.Seek(int64(position), 0)
	var remainingDocuments []map[string]interface{}
	for {
		var doc map[string]interface{}
		if err := dec.Decode(&doc); err != nil {
			break
		}
		remainingDocuments = append(remainingDocuments, doc)
	}

	// Move back to the position of the deleted document and truncate the file
	cube.File.Seek(int64(position), 0)
	cube.File.Truncate(int64(position))

	// Write back the remaining documents after the deleted position
	for _, doc := range remainingDocuments {
		if err := enc.Encode(doc); err != nil {
			return err
		}
	}

	return nil
}

func (cube *DocumentCube) calculatePosition(x, y, z int) int {
	return (z*cube.Size*cube.Size + y*cube.Size + x) * (int(unsafe.Sizeof(KeyValueBlock{})))
}

func main() {
	filePath := "document_cube.bin"
	defer os.Remove(filePath)

	documentCube, err := NewDocumentCube(filePath, cubeSize)
	if err != nil {
		fmt.Println("Error creating DocumentCube:", err)
		return
	}

	// Set documents with multiple key-value pairs
	document1 := map[string]interface{}{
		"name":   "John",
		"age":    25,
		"city":   "New York",
		"gender": "Male",
	}
	err = documentCube.SetDocument(0, 0, 0, document1)
	if err != nil {
		fmt.Println("Error setting document:", err)
		return
	}

	document2 := map[string]interface{}{
		"name":   "Jane",
		"age":    30,
		"city":   "San Francisco",
		"gender": "Female",
	}
	err = documentCube.SetDocument(1, 1, 1, document2)
	if err != nil {
		fmt.Println("Error setting document:", err)
		return
	}

	document3 := map[string]interface{}{
		"name":   "Bob",
		"age":    25,
		"city":   "New York",
		"gender": "Male",
	}
	err = documentCube.SetDocument(2, 2, 2, document3)
	if err != nil {
		fmt.Println("Error setting document:", err)
		return
	}

	// Print the cube contents before deletion
	fmt.Println("Cube Contents (Before Deletion):")
	for z := 0; z < cubeSize; z++ {
		for y := 0; y < cubeSize; y++ {
			for x := 0; x < cubeSize; x++ {
				document, err := documentCube.GetDocument(x, y, z)
				if err == nil {
					fmt.Printf("(%d, %d, %d): %v\n", x, y, z, document)
				}
			}
		}
	}

	// Delete a document
	err = documentCube.DeleteDocument(1, 1, 1)
	if err != nil {
		fmt.Println("Error deleting document:", err)
		return
	}

	// Print the cube contents after deletion
	fmt.Println("\nCube Contents (After Deletion):")
	for z := 0; z < cubeSize; z++ {
		for y := 0; y < cubeSize; y++ {
			for x := 0; x < cubeSize; x++ {
				document, err := documentCube.GetDocument(x, y, z)
				if err == nil {
					fmt.Printf("(%d, %d, %d): %v\n", x, y, z, document)
				}
			}
		}
	}
}
