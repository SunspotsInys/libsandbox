#include <stdio.h>
int main()
{
	int n,x,y;
	scanf("%d",&n);
	int i;
	for(i=0;i<n;i++){
		scanf("%d %d",&x,&y);
		printf("%d\n",x+y);
	}
}
