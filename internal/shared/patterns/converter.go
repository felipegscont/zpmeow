package patterns

// Converter defines the interface for converting between different types
type Converter[T, U any] interface {
	Convert(input T) U
}

// BatchConverter defines the interface for converting slices
type BatchConverter[T, U any] interface {
	ConvertBatch(input []T) []U
}

// BidirectionalConverter defines the interface for two-way conversion
type BidirectionalConverter[T, U any] interface {
	Convert(input T) U
	ConvertBack(input U) T
}

// EntityToDTOConverter defines the interface for entity to DTO conversion
type EntityToDTOConverter[Entity, DTO any] interface {
	ToDTO(entity Entity) DTO
	ToDTOBatch(entities []Entity) []DTO
}

// DTOToEntityConverter defines the interface for DTO to entity conversion
type DTOToEntityConverter[DTO, Entity any] interface {
	ToEntity(dto DTO) Entity
	ToEntityBatch(dtos []DTO) []Entity
}
