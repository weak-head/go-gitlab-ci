package models

// Gist is a declaration of a new source code gist.
//
//	@Description	Gist is a definition of a Source Code gist that should be saved to the storage.
type Gist struct {

	// Name is a human readable Gist name.
	Name string `json:"name" binding:"required" example:"Generate unique ID"`

	// Description is a human readable Gist description.
	Description string `json:"description" example:"Example of how to generate a unique ID in java script."`

	// Language is a programming language that is used in the gist.
	Language string `json:"language" binding:"required" example:"javascript"`

	// Code is a Source Code gits.
	Code string `json:"code" binding:"required" example:"for (let i = 0; i < 5; i++) {...}"`
}

// GistInfo provides a high-level information about the Gist.
//
//	@Description	GistInfo provides the descriptive information about the
//	@Description	Gist entry and doesn't return the actual Gist definition.
type GistInfo struct {
	// Id is a globally unique Gist ID that identifies this Gist entry.
	Id string `json:"id" binding:"required" example:"d17043a0-216c-4c56-9127-b0bf5e3a4c16"`

	// Name is a human readable Gist name.
	Name string `json:"name" binding:"required" example:"Generate unique ID"`

	// Description is a human readable Gist description.
	Description string `json:"description" example:"Example of how to generate a unique ID in java script."`

	// Language is a programming language that is used in the gist.
	Language string `json:"language" binding:"required" example:"javascript"`
}

// GistDetails provided the detailed information about Gist entry.
//
//	@Description	GistDetails provides the detailed information about the
//	@Description	requested Gist and includes all the available
//	@Description	public information that is stored in the service.
type GistDetails struct {
	// Id is a globally unique Gist ID that identifies this Gist entry.
	Id string `json:"id" binding:"required" example:"d17043a0-216c-4c56-9127-b0bf5e3a4c16"`

	// Name is a human readable Gist name.
	Name string `json:"name" binding:"required" example:"Generate unique ID"`

	// Description is a human readable Gist description.
	Description string `json:"description" example:"Example of how to generate a unique ID in java script."`

	// Language is a programming language that is used in the gist.
	Language string `json:"language" binding:"required" example:"javascript"`

	// Code is a Source Code gits.
	Code string `json:"code" binding:"required" example:"for (let i = 0; i < 5; i++) {...}"`

	// CreatedAt defines the date and time when the gist has been created.
	// This field uses RFC 3339 as the standard for the date-time format.
	CreatedAt string `json:"createdAt" binding:"required" example:"2023-06-07T18:27:25-04:00"`

	// LastUpdated defines the date and time when the gist has been updated.
	// This field uses RFC 3339 as the standard for the date-time format.
	LastUpdated string `json:"lastUpdated,omitempty" example:"2023-06-11T10:44:17-04:00"`

	// LastAccessed defines the date and time when the gist has been accessed.
	//
	// Only direct operations on this gist update the field.
	// Such operations as 'get all gists' doesn't update the field.
	//
	// This field uses RFC 3339 as the standard for the date-time format.
	LastAccessed string `json:"lastAccessed,omitempty" example:"2023-06-24T08:13:59-04:00"`
}
