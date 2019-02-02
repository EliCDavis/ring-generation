# Ring Generation

Attemping to generate meshes in obj format of rings to 3D print. End goal is to enscribe text onto the ring.

To run use:
```bash
docker-compose run ring-gen
```

## Current Progress

![Ring](https://i.imgur.com/pFQCKw4.png)

## Resources Used

### OBJ File Format and Material

 * [Obj file format](http://paulbourke.net/dataformats/obj/)
 * [More Obj file format](https://www.cs.cmu.edu/~mbz/personal/graphics/obj.html)
 * [Interfacing obj with mtl](https://people.cs.clemson.edu/~dhouse/courses/405/docs/brief-obj-file-format.html)
 * [Mtl file format](http://www.paulbourke.net/dataformats/mtl/)

### Triangulation

 * [Some Definitions](http://www.cs.cmu.edu/~quake/triangle.defs.html#cdt)
 * [More Definitions](http://www.cs.cmu.edu/~quake/triangle.delaunay.html)
 * [Unstructured Mesh Generation and Adaptation](https://hal.inria.fr/hal-01438967/document)
 * [MIT Unstructured Mesh Generation](https://popersson.github.io/pub/persson06unstructured.pdf)
 * [Berkley Lecture Notes on Delaunay Mesh Generation](https://people.eecs.berkeley.edu/~jrs/meshpapers/delnotes.pdf)
 * [Random Lecture Unstructured Mesh Generation](https://nptel.ac.in/courses/112106061/Module_3/Lecture_3.6.pdf)
 * [Recent Progress in Robust and Quality Delaunay Mesh Generation](https://pdfs.semanticscholar.org/fcbb/48fb3f7259dc8151adb65c43da605a121ca2.pdf)

### Other Random Geometry

 * [Sorting points clockwise order](https://stackoverflow.com/questions/6989100/sort-points-in-clockwise-order)
 * [Determining if a point lies within a polygon](https://www.geeksforgeeks.org/how-to-check-if-a-given-point-lies-inside-a-polygon/)
