#version 120
attribute vec3 obj_coord;
varying vec3 f_color;
uniform mat4 m_transform;

void main(void) {
	gl_Position = m_transform * vec4(obj_coord, 1.0);
	f_color = vec3(0.0,0.0,1.0);
}
