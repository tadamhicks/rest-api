FROM scratch
ADD rest-api /
EXPOSE 8000
CMD ["/rest-api"]