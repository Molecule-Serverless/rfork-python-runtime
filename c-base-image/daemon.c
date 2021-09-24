#include <stdio.h>
#include <string.h>
#define _GNU_SOURCE
#include <sched.h>
#include <sys/socket.h>
#include <sys/wait.h>
#include <sys/un.h>
#include <errno.h>
#include <unistd.h>
#include <assert.h>
#include <stdlib.h>

#include "syscall.h"

#define SOCKET_NAME "fork.sock"
#define RECEIVE_FD_COUNT 5
#define MAX_FD_COUNT 128
#define PID_BUF_LENGTH 32

// TODO: dump key should be passed by runc
#define DUMP_KEY 73

int swap_device_fd;
int fork_socket_fd;
int receive_fds(int fd, int fd_array[]);

int handle_fork_request(int fd) {
    // called by main process
    // will fork children and call os.swap
    int fd_array[5]; // 5 file descriptors to be received
    int ret;
    pid_t pid;
    int target_fd, uts_namespace_fd, pid_namespace_fd, ipc_namespace_fd, mnt_namespace_fd;

    ret = receive_fds(fd, fd_array);
    if (ret) {
        printf("receive_fds failed");
        return ret;
    }
    target_fd = fd_array[0];
    uts_namespace_fd = fd_array[1];
    pid_namespace_fd = fd_array[2];
    ipc_namespace_fd = fd_array[3];
    mnt_namespace_fd = fd_array[4];

    printf("before fork\n");

    pid = fork();
    if (pid) {
        // the grand parent process
        printf("in grand parent process\n");
        close(fd);
        close(target_fd);
        close(uts_namespace_fd);
        close(pid_namespace_fd);
        close(ipc_namespace_fd);
        close(mnt_namespace_fd);
        waitpid(pid, NULL, 0);
    } else {
        // the parent process
        printf("in parent process\n");
        fchdir(target_fd);
        chroot(".");
        close(target_fd);
        setns(uts_namespace_fd, 0);
        setns(pid_namespace_fd, 0);
        setns(ipc_namespace_fd, 0);
        setns(mnt_namespace_fd, 0);
        close(uts_namespace_fd);
        close(pid_namespace_fd);
        close(ipc_namespace_fd);
        close(mnt_namespace_fd);
        pid = fork();
        if (pid) {
            // the parent process
            char buf[PID_BUF_LENGTH];
            sprintf(buf, "%d", pid);
            size_t len = strlen(buf);
            printf("str length: %lu\n", len);
            size_t send_len = send(fd, buf, len, 0 /* flags */);
            printf("send length: %lu\n", send_len);
            exit(0);
        } else {
            // the child process
            printf("in client process\n");
            close(fork_socket_fd);

            // call swap here
            assert(swap_device_fd > 0);
            printf("swap_device_fd: %d\n", swap_device_fd);
            call_swap(swap_device_fd, DUMP_KEY);
            printf("should never reach here");
            perror("call_swap");
            assert(0);
        }
    }
    return 0;
}

int receive_fds(int fd, int fd_array[]) {
    struct msghdr msg;
    char message_buffer[1];
    struct cmsghdr* cmsg;
    int* data;
    int cmsg_length;
    int ret;
    int received_length = RECEIVE_FD_COUNT * sizeof(int);

    memset(&msg, 0, sizeof(msg));

    union {
        char buf[CMSG_SPACE(received_length)];
        struct cmsghdr align;
    } control_msg;

    memset(&control_msg, 0, sizeof(control_msg));

    struct iovec io = {
        .iov_base = message_buffer,
        .iov_len = sizeof(message_buffer)
    };
    
    msg.msg_name = NULL;
    msg.msg_namelen = 0;
    msg.msg_iov = &io;
    msg.msg_iovlen = 1;

    /*cmsghdr = (struct cmsghdr*)buf;
    cmsghdr->cmsg_len = CMSG_LEN(received_length);
    cmsghdr->cmsg_level = SOL_SOCKET;
    cmsghdr->cmsg_type = SCM_RIGHTS;*/

    msg.msg_control = control_msg.buf;
    msg.msg_controllen = CMSG_SPACE(received_length);

    if ((ret = recvmsg(fd, &msg, 0)) < 0) {
        printf("ret: %d, errno: %d\n", ret, errno);
        perror("recvmsg");
        return -1;
    }

    printf("recvmsg returns %d\n", ret);

    cmsg = CMSG_FIRSTHDR(&msg);
    data = (void*)CMSG_DATA(cmsg);
    cmsg_length = cmsg->cmsg_len;
    
    assert(cmsg_length - sizeof(struct cmsghdr) == received_length);
    for (int i = 0; i < RECEIVE_FD_COUNT; i++) {
        fd_array[i] = data[i];
        printf("recv fd: %d\n", data[i]);
    }

    printf("cmsg_length: %d\n", cmsg_length);
    return 0;
}

int main() {
    int fd;
    struct sockaddr_un addr;

    swap_device_fd = sopen();
    if (swap_device_fd < 0) {
        perror("sopen");
        assert(0);
    }

    // create socket fd
    if ((fd = socket(AF_UNIX, SOCK_STREAM, 0)) < 0) {
        perror("socket");
        return 0;
    }

    if (unlink(SOCKET_NAME)) {
        if (errno != ENOENT) {
            perror("unlink");
            return -1;
        }
    }

    memset(&addr, 0, sizeof(addr));
    addr.sun_family = AF_UNIX;
    strcpy(addr.sun_path, SOCKET_NAME);

    // bind the socket fd with a unix domain socket
    if (bind(fd, (struct sockaddr*)&addr, sizeof(addr))) {
        perror("bind");
        return -1;
    }

    if (listen(fd, 0)) {
        perror("listen");
        return -1;
    }

    fork_socket_fd = fd;

    // receive fds from unix domain socket
    while (1) {
        int accept_fd = accept(fd, NULL, NULL);
        if (accept_fd < 0) {
            perror("accept");
            return -1;
        }
        handle_fork_request(accept_fd);
        int ret = close(swap_device_fd);
        if (ret != 0) {
            perror("close");
            assert(0);
        }
        swap_device_fd = sopen();
        if (swap_device_fd < 0) {
            perror("sopen");
            assert(0);
        }
    }

    return 0;
}
