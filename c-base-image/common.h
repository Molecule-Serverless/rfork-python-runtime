#ifndef OS_SWAP_COMMON_H
#define OS_SWAP_COMMON_H

enum LibSwapCmd
{
  Nil = 0,  // for test only
  Dump = 1, // dump the memory mapping of this process
  // XD: currently, I don't know why we cannot set the cmd in `ioctrl` to be 2
  Swap = 3, // swap to another process
};

#endif