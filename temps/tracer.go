package temps

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
)

func TracerFrame() {
	// Open the JSON file
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer file.Close() // Defer closing the file until the function returns

	// Decode the JSON content into the data structure
	var data Data
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	for i := 0; i < len(data.Models); i++ {
		data.Models[i].LowerName = strings.ToLower(data.Models[i].Name)
		data.Models[i].AppName = data.AppName
		data.Models[i].ProjectName = data.ProjectName
		rl_list := make([]Relationship, 0)
		for k := 0; k < len(data.Models[i].RlnModel); k++ {
			rmf := strings.Split(data.Models[i].RlnModel[k], "$")
			cur_relation := Relationship{
				ParentName:      data.Models[i].Name,
				LowerParentName: data.Models[i].LowerName,
				FieldName:       rmf[0],
				LowerFieldName:  strings.ToLower(rmf[0]),
				MtM:             rmf[1] == "mtm",
				OtM:             rmf[1] == "otm",
				MtO:             rmf[1] == "mto",
			}
			rl_list = append(rl_list, cur_relation)
			data.Models[i].Relations = rl_list
		}

		for j := 0; j < len(data.Models[i].Fields); j++ {
			data.Models[i].Fields[j].BackTick = "`"
			cf := strings.Split(data.Models[i].Fields[j].CurdFlag, "$")

			data.Models[i].Fields[j].Get, _ = strconv.ParseBool(cf[0])
			data.Models[i].Fields[j].Post, _ = strconv.ParseBool(cf[1])
			data.Models[i].Fields[j].Patch, _ = strconv.ParseBool(cf[2])
			data.Models[i].Fields[j].Put, _ = strconv.ParseBool(cf[3])
			data.Models[i].Fields[j].AppName = data.AppName
			data.Models[i].Fields[j].ProjectName = data.ProjectName

		}
	}

	// ############################################################
	common_tmpl, err := template.New("data").Parse(tracerTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("observe", os.ModePerm)
	if err != nil {
		panic(err)
	}

	common_file, err := os.Create("observe/tracer.go")
	if err != nil {
		panic(err)
	}
	defer common_file.Close()

	err = common_tmpl.Execute(common_file, data)
	if err != nil {
		panic(err)
	}

	// running go mod tidy finally
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		fmt.Printf("error: %v \n", err)
	}

}

var tracerTemplate = `
package observe

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"{{.ProjectName}}.com/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var AppTracer = otel.Tracer(fmt.Sprintf("fiber-server-%v", configs.AppConfig.Get("APP_NAME")))

func InitTracer() *sdktrace.TracerProvider {
	traceExporter := configs.AppConfig.Get("TRACE_EXPORTER")
	tracerHost := configs.AppConfig.Get("TRACER_HOST")
	tracerPort := configs.AppConfig.GetOrDefault("TRACER_PORT", "9411")

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(configs.AppConfig.Get("APP_NAME")),
		)),
	)
	// app logger with jager

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	// otel.SetErrorHandler(&otelErrorHandler{logger: app_logger})

	const traceExporterFiber = "fiber"

	if (traceExporter != "" && tracerHost != "") || traceExporter == traceExporterFiber {
		var (
			exporter sdktrace.SpanExporter
			// err      error
		)

		switch strings.ToLower(traceExporter) {
		case "jaeger":
			// app_logger.Log("Exporting traces to jaeger.")

			exporter, _ = otlptracegrpc.New(context.Background(), otlptracegrpc.WithInsecure(),
				otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%s", tracerHost, tracerPort)))

			batcher := sdktrace.NewBatchSpanProcessor(exporter)
			tp.RegisterSpanProcessor(batcher)
		}
	}
	return tp
}

func AppSpanner(ctx context.Context, span_name string) (context.Context, oteltrace.Span)  {
	gen, _ := uuid.NewV7()
	id := gen.String()
	trace, span := AppTracer.Start(ctx, span_name, oteltrace.WithAttributes(attribute.String("id", id)))
	return trace, span
}



`
