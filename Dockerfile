FROM scratch
COPY agriTracking /
COPY node_modules /node_modules
COPY templates /templates
COPY authz_model.conf /
COPY authz_policy.csv /
WORKDIR /
CMD ["/agriTracking"]
