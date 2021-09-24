#include <assert.h>
#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <time.h>

#include "syscall.h"

int a = 0;
const int b = 1;

#if 0
void print_fs_28();

static inline long get_fs_28() {
    long ret = 0;
    __asm__ __volatile__ ("mov %%fs:0x28, %0":"=r"(ret));
    return ret;
}

void print_fs_28() {
    long ret = get_fs_28();
    printf("fs:0x28: 0x%lx\n", ret);
    return;
}
#endif

int
main()
{
  struct timespec spec;
  clock_gettime(CLOCK_REALTIME, &spec);
  fprintf(stderr, "%ld %ld\n", spec.tv_sec, spec.tv_nsec); // checkpoint 1: start

  int sd = sopen();
  if (sd <= 0) {
      printf("error opening device");
      perror("sopen");
      assert(0);
  }

  clock_gettime(CLOCK_REALTIME, &spec);
  fprintf(stderr, "%ld %ld\n", spec.tv_sec, spec.tv_nsec); // checkpoint 2: after open device

  int res = call_swap(sd, 73);

  printf("call swap res: %d\n", res);

  assert(0);

  while (1) {
      sleep(1);
  }

  if (b > 0) {

    while (1) {
      // should never execute here
      sleep(1);
    }
  }

  return 0;
}
