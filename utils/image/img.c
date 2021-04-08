#include "img.h"


int sum(int a, int b){
    return (a+b);
}

void address(char * path){
    printf("From C: %s\n", path);
}

int *makeArr(int size, char* str){
    int *arr = (int *)malloc(size*sizeof(int));
    srand((unsigned)time(NULL));

    printf("%s\n", str);
    for (int i=0;i<size;i++){
        arr[i] = rand()%100;
    }
    return arr;
}
