package component

import "pkg.world.dev/world-engine/cardinal"

func IsCollisonThere(world cardinal.WorldContext, pos Position) (bool, error) {
	found, id, err := pos.GetEntityIDByPosition(world)
	if err != nil {
		return false, err
	}
	if found {
		colType, err := cardinal.GetComponent[Collidable](world, id)
		if err != nil {
			return false, err
		}
		switch colType.Type {
		case ItemCollide:
			// ok to overlap items
			return false, nil
		default:
			// not ok to overlap other types of collidable
			return true, nil
		}
	}
	// no entity found, so it's not a wall
	return false, nil
}
