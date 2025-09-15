package patterns

type Converter[T, U any] interface {
	Convert(input T) U
}

type BatchConverter[T, U any] interface {
	ConvertBatch(input []T) []U
}

type BidirectionalConverter[T, U any] interface {
	Convert(input T) U
	ConvertBack(input U) T
}

type EntityToDTOConverter[Entity, DTO any] interface {
	ToDTO(entity Entity) DTO
	ToDTOBatch(entities []Entity) []DTO
}

type DTOToEntityConverter[DTO, Entity any] interface {
	ToEntity(dto DTO) Entity
	ToEntityBatch(dtos []DTO) []Entity
}
