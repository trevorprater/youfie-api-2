import os

for f in os.listdir('../routers'):
    routes = []
    methods = []
    with open('../routers/' + f) as ff:
        route = methods = ''
        for line in ff.readlines():
            if 'router.Handle(' in line or 'router.HandleFunc(' in line:
                route = line[line.find('(') + 1:line.find(',')].replace(
                    '"', '').strip()
            if 'Methods(' in line:
                methods = line.split('Methods(')[-1].replace(')', '').replace(
                    '"', '').strip()
            if not route == '' and not methods == '':
                print 'endpoint: {},  methods: {}'.format(route, methods)
                route = methods = ''
