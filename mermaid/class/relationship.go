package class

// Relationship is a type that represents the relationship between classes.
//
// Ref. https://mermaid.js.org/syntax/classDiagram.html#defining-relationship
type Relationship string

const (
	// RelationshipInheritance is inheritance from right class to left class.
	// e.g. Vehicle <|-- Car
	RelationshipInheritance Relationship = "<|--"
	// RelationshipInheritanceReverse is inheritance from left class to right class.
	// e.g. Car --|> Vehicle
	RelationshipInheritanceReverse Relationship = "--|>"
	// RelationshipComposition is composition from left class to right class.
	// e.g. House *-- Room
	RelationshipComposition Relationship = "*--"
	// RelationshipCompositionReverse is composition from right class to left class.
	// e.g. Room --* House
	RelationshipCompositionReverse Relationship = "--*"
	// RelationshipAggregation is aggregation from left class to right class.
	// e.g. Team o-- Player
	RelationshipAggregation Relationship = "o--"
	// RelationshipAggregationReverse is aggregation from right class to left class.
	// e.g. Player --o Team
	RelationshipAggregationReverse Relationship = "--o"
	// RelationshipAssociation is association from left class to right class.
	// e.g. Order --> Customer
	RelationshipAssociation Relationship = "-->"
	// RelationshipAssociationReverse is association from right class to left class.
	// e.g. Customer <-- Order
	RelationshipAssociationReverse Relationship = "<--"
	// RelationshipLink is a solid link.
	// e.g. A -- B
	RelationshipLink Relationship = "--"
	// RelationshipDependency is dependency from left class to right class.
	// e.g. Service ..> Repository
	RelationshipDependency Relationship = "..>"
	// RelationshipDependencyReverse is dependency from right class to left class.
	// e.g. Repository <.. Service
	RelationshipDependencyReverse Relationship = "<.."
	// RelationshipRealization is realization from left class to right class.
	// e.g. Driver ..|> Drivable
	RelationshipRealization Relationship = "..|>"
	// RelationshipRealizationReverse is realization from right class to left class.
	// e.g. Drivable <|.. Driver
	RelationshipRealizationReverse Relationship = "<|.."
	// RelationshipDashedLink is a dashed link.
	// e.g. A .. B
	RelationshipDashedLink Relationship = ".."
)
