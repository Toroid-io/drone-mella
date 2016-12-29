FROM toroid/mella
ADD drone-mella /bin/
ENTRYPOINT /bin/drone-mella
