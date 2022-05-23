{{$svrName := .ServiceName}}
type {{.ServiceName}}GinServer interface {
{{- range .MethodSets}}
    {{- .Comment}}
    {{- if (.HasFile)}}
    {{.MethodName}}(context.Context, *multipart.FileHeader) (*{{.Reply}}, error)
    {{- else}}
	{{.MethodName}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
	{{- end}}
{{- end}}
}

type {{.ServiceName}} struct{
	server {{.ServiceName}}GinServer
	router gin.IRouter
}

func Register{{.ServiceName}}(r gin.IRouter, srv {{.ServiceName}}GinServer) {
	s := {{.ServiceName}}{
		server: srv,
		router:     r,
	}
	s.RegisterHandlers()
}

func response(ctx *gin.Context, status, code int, msg string, data interface{}) {
	ctx.JSON(status, map[string]interface{}{
		"code": code,
		"msg": msg,
		"data": data,
	})
}

{{range .MethodSets}}
func (s *{{$svrName}}) {{.MethodName}} (ctx *gin.Context) {
{{- if not (.HasFile)}}
	var in {{.Request}}
{{- end}}
{{- if .HasPathParams}}
	if err := ctx.ShouldBindUri(&in); err != nil {
		response(ctx, 400, -1, "参数错误", nil)
		return
	}
{{- else if .HasByte}}
    data, err := ioutil.ReadAll(ctx.Request.Body)
    	if err != nil {
    		response(ctx, 400, -1, "参数错误", nil)
    		return
    	}
	in.HttpBody = &httpbody.HttpBody{
		ContentType: ctx.ContentType(),
		Data:        data,
	}
{{- else if .HasFile}}
    var in *multipart.FileHeader
    in, err := ctx.FormFile("file")
    if err != nil {
        response(ctx, 400, -1, "参数错误", nil)
        return
    }
{{- else if eq .Method "GET" "DELETE" }}
	if err := ctx.ShouldBindQuery(&in); err != nil {
		response(ctx, 400, -1, "参数错误", nil)
		return
	}
{{- else if eq .Method "POST" "PUT" }}
	if err := ctx.ShouldBindJSON(&in); err != nil {
		response(ctx, 400, -1, "参数错误", nil)
		return
	}
{{- else}}
	if err := ctx.ShouldBind(&in); err != nil {
		response(ctx, 400, -1, "参数错误", nil)
		return
	}
{{- end}}
	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	{{- if not (.HasFile)}}
    	out, err := s.server.{{.MethodName}}(newCtx, &in)
    {{- else}}
        out, err := s.server.{{.MethodName}}(newCtx, in)
    {{- end}}
	if err != nil {
		_ = ctx.Error(err)
		response(ctx, 500, -1, "未知错误", nil)
		return
	}

	response(ctx, 200, 0, "成功", out)
}
{{end}}

func (s *{{$svrName}}) RegisterHandlers() {
{{- range .MethodSets}}
    {{- .Comment}}
	s.router.{{.Method}}("{{.Path}}", s.{{.MethodName}})
{{- end}}
}
