package dbqueries

import "github.com/jackc/pgx/v5"

const PostAlbumByProjectId = `
	INSERT INTO album (project_id, name, cover)
	VALUES
	(@projectId, @name, @cover);
`

func PostAlbumByProjectIdArgs(projectId, name, cover string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"projectId": projectId,
		"name":      name,
		"cover":     cover,
	}
}

const DeleteAlbumById = `
	DELETE FROM album
	WHERE
	album_id = @albumId
`

func DeleteAlbumByIdArgs(albumId string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"albumId": albumId,
	}
}

const PatchAlbumMetadataById = `
	UPDATE album
	SET name = @name
	WHERE
	album_id = @albumId
`

func PatchAlbumMetadataByIdArgs(albumId, name string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"albumId": albumId,
		"name":    name,
	}
}

const PatchAlbumCoverById = `
	UPDATE album
	SET cover = @cover
	WHERE
	album_id = @albumId
`

func PatchAlbumCoverByIdArgs(albumId, cover string) pgx.NamedArgs {
	return pgx.NamedArgs{
		"albumId": albumId,
		"cover":   cover,
	}
}