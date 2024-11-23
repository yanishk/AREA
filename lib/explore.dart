import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "custom_widget.dart";
import "footer_page.dart";

class ExplorePage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {

    return Scaffold(
      body: SingleChildScrollView (
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            TopBar(),
            MyDivider(),
            SearchBar(),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 1),
    );
  }
}

class SearchBar extends StatelessWidget {
  const SearchBar({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: Container(
        decoration: BoxDecoration(
          color: const Color.fromARGB(255, 238, 238, 238),
          borderRadius:  BorderRadius.circular(3),
        ),
        child: TextField(
          decoration: InputDecoration(
            icon: Icon(Icons.search, size: 35,),
            iconColor: Colors.black,

            hintStyle: TextStyle(fontSize: 20, color: const Color.fromARGB(255, 156, 156, 156), fontWeight: FontWeight.bold,),
            hintText: 'Search your trips',
            border: InputBorder.none,
            contentPadding: EdgeInsets.symmetric(vertical: 10, horizontal: 10),
          ),
        ),
      ),
    );
  }
}

class TopBar extends StatelessWidget {
  const TopBar({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Row (
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        TextButton(
          onPressed: () {},
          child: AutoSizeText("All", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black), maxLines: 1,),
        ),

        TextButton(
          onPressed: () {},
          child: AutoSizeText("Applets", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black), maxLines: 1,),
        ),

        TextButton(
          onPressed: () {},
          child: AutoSizeText("Services", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black), maxLines: 1,),
        ),

        TextButton(
          onPressed: () {},
          child: AutoSizeText("Stories", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black), maxLines: 1,),
        ),
      ],
    );
  }
}
