#version 120
attribute vec3 obj_coord;
attribute vec3 obj_normal;
uniform mat4 m_transform;
varying vec3 f_normal;
void main(void) {
	gl_Position = m_transform * vec4(obj_coord, 1.0);
	f_normal = obj_normal;
}
