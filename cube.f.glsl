#version 120
varying vec3 f_normal;
void main(void) {
	gl_FragColor[0] = f_normal.x;
	gl_FragColor[1] = f_normal.y;
	gl_FragColor[2] = f_normal.z;
}
