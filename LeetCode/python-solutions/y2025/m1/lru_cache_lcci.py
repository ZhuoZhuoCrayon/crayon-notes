import collections


class LRUCache:
    """# https://leetcode.cn/problems/lru-cache-lcci/description/"""
    
    def __init__(self, capacity: int):
        self._cache = collections.OrderedDict()
        self._capacity: int = capacity

    def get(self, key: int):
        value = self._cache.get(key)
        if value is None:
            return -1

        self._cache.move_to_end(key)
        return value

    def put(self, key, value) -> None:
        if len(self._cache) >= self._capacity and key not in self._cache:
            self._cache.popitem(last=False)
        self._cache[key] = value
        self._cache.move_to_end(key)


if __name__ == '__main__':
    cache = LRUCache(2)
    print(cache.get(2))
    cache.put(2, 6)
    print(cache.get(1))

    cache.put(1, 5)
    cache.put(1, 2)

    print(cache.get(1))
    print(cache.get(2))

