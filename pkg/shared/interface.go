package shared

type Mapper[From any, To any] interface {
	FromModel(From) To
	ToModel() From
}
