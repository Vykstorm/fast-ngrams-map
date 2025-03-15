package src

// The following interface represents a mapping between strings and and type of values
type Map[T any] interface {
	/**
	Returns the value associated with the given key and true as the second return value.
	if there is an entry with the given key in the map. If the key
	is not in the map, it returns null for type T and false.

	This method is not thread-safe. Put() and Get() can't be called at the same time.
	However, multiple reads (calls to Get(), Has() or Len() ) can be safely performed concurrently.

	The key must have non-zero length
	*/
	Get(key string) (T, bool)

	/**
	Returns true if there is an entry with the given key in the map. False otherwise.
	This method is not thread-safe. Put() and Has() can't be called at the same time.
	However, multiple reads (calls to Get(), Has() or Len() ) can be safely performed concurrently.

	The key must have non-zero length
	*/
	Has(key string) bool

	/**
	Puts a new value in the map with the given key.
	This method is not thread-safe. Put() and Get() can't be called at the same time.
	However, multiple reads (calls to Get(), Has() or Len() ) can be safely performed concurrently.
	The key must have non-zero length
	*/
	Put(key string, value T)

	/**
	Returns the number of entries in this map.
	This method is not thread-safe. Put() and Len() can't be called at the same time.
	However, multiple reads (calls to Get(), Has() or Len() ) can be safely performed concurrently.
	*/
	Len() int
}
