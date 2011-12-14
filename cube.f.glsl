#version 120
varying vec3 vary_norm;

struct lightsource {
	vec4 position;
	vec4 diffuse;
};
lightsource light0 = lightsource(
	vec4(1.0, 1.0, 1.0, 0.0),
	vec4(1.0, 1.0, 1.0, 1.0)
);

struct material {
	vec4 diffuse;
};
material monkey = material(vec4(1.0, 0.8, 0.8, 1.0));

void main(void) {

	vec3 normal_dir = normalize(vary_norm);
	vec3 light_dir = normalize(vec3(light0.position));
	
	float d_r = max(0.0, dot(normal_dir, light_dir));
	gl_FragColor = vec4(vec3(light0.diffuse) * vec3(monkey.diffuse) * d_r, 1.0);
}
