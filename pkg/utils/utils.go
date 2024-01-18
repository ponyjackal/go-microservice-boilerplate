package utils

import (
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ponyjackal/go-microservice-boilerplate/pkg/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// basepath is the root directory of this package.
var basepath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

// Path returns the absolute path the given relative file or directory path,
// relative to the directory in the
// user's GOPATH.  If rel is already absolute, it is returned unmodified.
func Path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basepath, rel)
}

// Utility function to convert time.Time to *timestamppb.Timestamp
func convertToTimestamp(t time.Time) string {
	if t.IsZero() {
		return "2023-01-01T00:00:00Z"
	}
	return t.Format(time.RFC3339)
}

func SplitByLastPeriod(value string) string {
	// Find the last occurrence of "."
	lastIndex := strings.LastIndex(value, ".")
	if lastIndex == -1 {
		return ""
	}
	// Split the string into two parts based on the last "."
	result := value[:lastIndex]
	return result
}

func ConvertProtoToByte(obj protoreflect.ProtoMessage) ([]byte, error) {
	marshaller := protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
	jsonString, err := marshaller.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return jsonString, nil
}

// ParseStringToTime converts a string to time.Time. Returns a default time if parsing fails.
func ParseStringToTime(ts string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		logger.Infof("Error parsing time: %s", err)
		return time.Unix(0, 0)
	}
	return parsedTime
}

func GRPCErrorHandler(c *gin.Context, err error) {
	if err != nil {
		grpcStatus, _ := status.FromError(err)
		switch grpcStatus.Code() {
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{"error": grpcStatus.Message()})
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": grpcStatus.Message()})
		case codes.PermissionDenied:
			c.JSON(http.StatusForbidden, gin.H{"error": grpcStatus.Message()})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{"error": grpcStatus.Message()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		}
		return
	}
}

func Contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
