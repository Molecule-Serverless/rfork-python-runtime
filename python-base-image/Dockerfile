FROM val01:5000/python-base-image:dev-base

COPY setup.py /

COPY ol.c /

RUN cd / && python3 setup.py build_ext --inplace && mv ol.*.so /ol.so

COPY daemon.py /
