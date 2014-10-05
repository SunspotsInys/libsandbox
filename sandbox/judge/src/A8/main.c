#include <stdio.h>
int main()
{
	int a, b,c;
	do{
		a=2;
		printf("请输入一个数字：");
		scanf("%d",&c);
		printf("%d\n",a);
		while(a<c)
		{
			a++;
			b=2;
			while( a>b)
			{
				if(a%b==0)
					break;
				b++;
			}

			if(a==b)
			{printf("%d\n",a);}
		}
	}
	while(c!=0);
	//system("pause");
	return 0;
}
