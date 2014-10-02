#include <stdio.h>
void p(int arr[],int m);
int main(){
	int i,n,m;
	scanf("%d",&n);
	for(i=0;i<n;i++){
		scanf("%d",&m);
		int arr [m+1];
		p(arr,m);
	}
}
void p(int arr[],int m){
	int i;
	for(i=1;i<m+1;i++){
		arr[i]=0;
	}
	for(i=2;i<m+1;i++){
		if(!arr[i]){
			int j ;
			for(j=2*i;j<m+1;j+=i){
				arr[j]=1;
			}
		}
	}
	int tail;
	for(i=m;i>>0;i--){
		if(!arr[i]){
			tail = i;
			break;
		}
	}
	for(i=2;i<m+1;i++){
		if(!arr[i]){
			if(i==tail)
				printf("%d ",i);
			else
				printf("%d ",i);
		}
	}
}
