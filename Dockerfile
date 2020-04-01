FROM golang:onbuild
# Add Maintainer Info
LABEL maintainer="ABHISHEK AGARWAL"

EXPOSE 8081
# Command to run the executable
CMD ["./main","0.0.0.0"]