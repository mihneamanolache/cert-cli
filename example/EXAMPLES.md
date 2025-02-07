# How To Use `cert-cli` For Data Analytics

## Directory Structure
The `example/` directory holds the following files:
- `targets.txt`: text file containing queries (domain/organization names) for which we want to gather information
- `collect.sh`: bash script for running `cert-cli` in multi-thread and speed up data gathering process
- `data/`: directory where we save the generated `.jsos` files
- `analyse.ipynb`: Jupyter notebook for data analysis

## 1. Collecting Certificate Data
### Step 1: Prepare Your Target List
Create a targets.txt file inside the example/ directory, listing domain names or organization names you want to analyze:
```bash
cat example/targets.txt
```

### Step 2: Run cert-cli to Gather Data
Use the provided `example/collect.sh` script to collect certificate information for multiple targets efficiently:
```bash
bash example/collect.sh
```
This script runs `cert-cli` in parallel for domains in `targets.txt` and saves results as JSON files inside `example/$DATE/`.

Alternatively, you can run cert-cli manually:
```bash
query=Dreamworks
cert-cli --q "$query" --match "ILIKE" --o "$query"
```

## 2. Analyzing the Collected Data
### Step 3: Open Jupyter Notebook
Open `example/analyse.ipynb` and run the provided PySpark-based analysis.

### Step 4: Process and Classify Domains
Inside the Jupyter Notebook:
- Load JSON data into a PySpark DataFrame:
```py
from pyspark.sql import SparkSession
spark = SparkSession.builder.appName("JSON Analysis").getOrCreate()
df = spark.read.option("multiline", "true").json("example/data/*.json").cache()
df.show()
```
- Extract domain details and classify them:
```py
from pyspark.sql.functions import col, explode, collect_set, size, desc, when, expr
out = df.select(explode(col("certificates")).alias("certificate"))
        .select(col("certificate.commonName").alias("domain"), col("certificate.san"))
        .select(col("domain"), explode(col("san")).alias("alternate_domain"))
out.show()
```
- Identify subdomains and external domains:
```py
classified = (
    out.withColumn("is_subdomain", col("alternate_domain").startswith(col("domain")))
    .groupBy("domain")
    .agg(
        collect_set(when(col("is_subdomain"), col("alternate_domain"))).alias("subdomains"),
        collect_set(when(~col("is_subdomain"), col("alternate_domain"))).alias("external_domains")
    )
)
classified.show(truncate=False)
```

# Conclusion
Using `cert-cli`, you can efficiently map an organization's web infrastructure and analyze its certificates using PySpark and Jupyter Notebook. This workflow enables security researchers to detect subdomains, external connections, and more.

For advanced usage, modify `analyse.ipynb` to include additional transformations, visualizations, or integrations with security monitoring tools.

Happy hacking! 