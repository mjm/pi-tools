package nomadic

func Logging(tag string) map[string]interface{} {
	return map[string]interface{}{
		"type": "journald",
		"config": []map[string]interface{}{
			{
				"tag": tag,
			},
		},
	}
}
