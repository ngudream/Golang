��ƽ̨���룺
�����ҵ��ǹ�����������������Go��Դ�������������б��밲װ����ɺ�Go/src�µ�make.bash���ɿ�ƽ̨�ı�������������ÿ����ʱ��Ҫָ���ǣϣϣӼ��ǣϣ��ңã������б��룬�磺
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build����������Ϊlinux 64λϵͳ�µĳ���
�����Ҫ����windows 64 λ�µĳ�����Ҫ��дһ�����
CGO_ENABLED=0 GOOS=windows GOARCH=amd64(��386) go build����������Ϊwindows 64λϵͳ�µĳ���

ͼ�ν��棺
ͨ�����main.manifestִ�У����walk��װ���ˣ�Ӧ�þ�����������
rsrc -manifest main.manifest -o rsrc.syso

��ͼ��rsrc����Ĳ����б�������Ϊ���ѵĳ���Ӹ�icoͼ�꣺
rsrc -manifest main.manifest �Cico icon.ico -o rsrc.syso

����ʾdos:
go build -ldflags "-H windowsgui"