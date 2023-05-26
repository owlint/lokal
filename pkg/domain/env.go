package domain

type EnvironmentVariable struct {
	Name  string
	Value string
}

// func MakeEnvironmentVariables(configEnvs []ConfigEnvironmentVariable, secrets []Secret) ([]EnvironmentVariable, error) {
// 	envs := []EnvironmentVariable{}

// 	for _, env := range configEnvs {
// 		if env.Value != "" {
// 			envs = append(envs, EnvironmentVariable{
// 				Name:  env.Name,
// 				Value: env.Value,
// 			})
// 		} else if env.Secret != "" {
// 			secret, err := getSecret(env.Secret, secrets)
// 			if err != nil {
// 				return nil, err
// 			}

// 			envs = append(envs, EnvironmentVariable{
// 				Name:  env.Name,
// 				Value: secret.Value,
// 			})
// 		} else {
// 			return nil, fmt.Errorf("environnement variable must have a value or reference a secret")
// 		}
// 	}
// 	return envs, nil
// }

// TODO
// func getSecret(name string, secrets []Secret) (*Secret, error) {
// 	for _, secret := range secrets {
// 		if secret.Name == name {
// 			return &secret, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("couldn't find secret '%s'", name)
// }
