#version 120
uniform mat4 m_inv_view;
varying vec3 vary_norm;
varying vec4 vary_pos;

struct lightSource
{
  vec4 position;
  vec4 diffuse;
  vec4 specular;
  float constantAttenuation, linearAttenuation, quadraticAttenuation;
  float spotCutoff, spotExponent;
  vec3 spotDirection;
};

//lightsource light0 = lightsource(
//	vec4(1.0, 1.0, 1.0, 0.0),
//	vec4(1.0, 1.0, 1.0, 1.0),
//
//);
lightSource light0 = lightSource(
  vec4(1.0,  2.0,  6.0, 1.0),
  vec4(1.0,  1.0,  1.0, 1.0),
  vec4(1.0,  1.0,  1.0, 1.0),
  0.0, 0.2, 0.1,
  180.0, 0.0,
  vec3(0.0, 0.0, 0.0)
);

struct material {
	vec4 ambient;
	vec4 diffuse;
	vec4 specular;
	float sheen;
};
material monkey = material(
	vec4(1.0, 1.0, 1.0, 1.0), //Ambient
	vec4(1.0, 0.8, 0.8, 1.0), //Diffuse
	vec4(1.0, 1.0, 1.0, 1.0), //Specular
	20.0
	);


vec4 latent_ambient = vec4(0.1, 0.1, 0.1, 1.0);

void main(void) {

	vec3 normal_dir = normalize(vary_norm);
	vec3 viewer_dir = normalize(vec3(m_inv_view * vec4(0.0, 0.0, 0.0, 1.0) - vary_pos));
	vec3 light_dir;
	float attenuation;



	//LOTS OF THINGS

	//Test directional
	if(light0.position.w == 0.0) {
		attenuation = 1.0;
		light_dir = normalize(vec3(light0.position));
	} else {
		vec3 offset_to_light = vec3(light0.position - vary_pos);
		float light_distance = length(offset_to_light);
		light_dir = normalize(offset_to_light);

		//Falloff
		attenuation = 1.0 / (
				light0.constantAttenuation +
				light0.linearAttenuation * light_distance +
				light0.quadraticAttenuation * light_distance * light_distance

			);

		//Spotlight cone
		if (light0.spotCutoff <= 90.0) {
			float clampedCos = max(0.0, dot(-light_dir, light0.spotDirection));

			if(clampedCos < cos(radians(light0.spotCutoff))) {
				attenuation = 0.0;
			} else {
				attenuation = attenuation * pow(clampedCos, light0.spotExponent);
			}
		}
	}


	vec3 ambient_portion = vec3(latent_ambient) * vec3(monkey.ambient);

	vec3 diffuse_portion = attenuation * vec3(light0.diffuse) * vec3(monkey.diffuse) * max(0.0, dot(normal_dir, light_dir));

	vec3 specular_portion = vec3(0.0, 0.0, 0.0);
	if(dot(normal_dir, light_dir) >= 0.0) { //On the right side
		specular_portion = attenuation * vec3(light0.specular) * vec3(monkey.specular) * 
		pow(max(0.0, dot(reflect(-light_dir, normal_dir), viewer_dir)), monkey.sheen);
	}

	gl_FragColor = vec4(ambient_portion + diffuse_portion + specular_portion, 1.0);
}
