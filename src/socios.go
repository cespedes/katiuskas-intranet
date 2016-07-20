package main

func sociosHandler(ctx *Context) {
	log(ctx, LOG_DEBUG, "Page /socios")

	var activos, baja_temporal int
	p := make(map[string]interface{})

	db.QueryRow("SELECT count(*) FROM vperson WHERE type >= 4").Scan(&activos)
	db.QueryRow("SELECT count(*) FROM vperson WHERE type = 3").Scan(&baja_temporal)
	p["activos"] = activos
	p["baja_temporal"] = baja_temporal

	renderTemplate(ctx, "socios", p)
}
